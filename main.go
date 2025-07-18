package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/PuerkitoBio/goquery"
)

var progressBar *widget.ProgressBar
var downloadButton *widget.Button

func main() {
	a := app.NewWithID("com.educabiz.downloader")
	w := a.NewWindow("Educabiz Downloader")
	w.Resize(fyne.NewSize(600, 400))

	instructionsText := "Instructions:\n" +
		"1. Fill in the subdomain\n" +
		"2. Press \"check subdomain\" to validate that it is correct\n" +
		"3. Fill in the username and password\n" +
		"4. Press download and wait for the progress bar to finish"
	instructions := widget.NewLabel(instructionsText)

	subdomainInput := widget.NewEntry()
	subdomainInput.SetPlaceHolder("subdomain goes here")
	subdomainInput.OnChanged = func(string) {
		downloadButton.Disable()
	}

	subdomainInputCheckButton := widget.NewButton("check subdomain", func() {
		// Query the school's homepage
		hostname := fmt.Sprintf("%s.educabiz.com", subdomainInput.Text)
		resp, err := http.Get(fmt.Sprintf("https://%s/", hostname))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// If the school does not exist the server redirects the request to a school not found page
		found, _ := regexp.MatchString(`/school/notFound`, resp.Request.URL.String())
		if found {
			fmt.Println("No such school")
			return
		} else {
			downloadButton.Enable()
			return
		}
	})
	subdomainLayout := container.New(layout.NewGridLayout(2), subdomainInput, subdomainInputCheckButton)

	usernameInput := widget.NewEntry()
	usernameInput.SetPlaceHolder("username")
	usernameInput.Validator = validation.NewRegexp(".+", "can't be empty")
	passwordInput := widget.NewPasswordEntry()
	passwordInput.SetPlaceHolder("password")
	passwordInput.Validator = validation.NewRegexp(".+", "can't be empty")
	credentialsLayout := container.New(layout.NewGridLayout(2), usernameInput, passwordInput)

	// folderSelectionDialog := dialog.NewFolderOpen(func(fyne.ListableURI, error) {}, w)
	// folderSelectionButton := widget.NewButton("Select Output Folder", func() {
	// 	folderSelectionDialog.Show()
	// })
	// folderSelectionButton.Disable()

	downloadButton = widget.NewButton("Download!", func() {
		go scrape(fmt.Sprintf("https://%s.educabiz.com/", subdomainInput.Text), usernameInput.Text, passwordInput.Text)
	})
	downloadButton.Disable()

	progressBar = widget.NewProgressBar()

	w.SetContent(container.NewVBox(
		instructions,
		subdomainLayout,
		credentialsLayout,
		// folderSelectionButton,
		downloadButton,
		progressBar,
	))

	w.ShowAndRun()
}

func scrape(hostname string, username string, password string) {
	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar: jar,
	}

	// Get a session cookie
	_, err := client.Get(hostname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	authenticityToken, _ := getAuthenticityToken(hostname, jar)

	// Check if credentials are valid
	resp, err := client.PostForm(
		fmt.Sprintf("%spublic/authenticateeducabiz", hostname),
		url.Values{
			"username": {username},
			"password": {password},
		},
	)
	fmt.Println(resp.StatusCode)
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	var result EBizAuthenticate
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	if result.Error != "" {
		fmt.Println("Failed to login")
		return
	}

	// Get the actual login cookies
	client.PostForm(
		fmt.Sprintf("%sauthenticate", hostname),
		url.Values{
			"authenticityToken": {authenticityToken},
			"username":          {username},
			"password":          {password},
		},
	)

	// Loop Children
	children := getChildrenList(client, hostname)
	fmt.Println(children)

	for _, child := range children {
		downloadDir := fmt.Sprintf("%s/gallery", child.Name)
		os.MkdirAll(downloadDir, 0750)
		pictures := getChildPhotos(client, hostname, child)
		totalPictures := len(pictures)
		for i, picture := range pictures {
			filePath := fmt.Sprintf("%s/%s-%d", downloadDir, picture.ShortDate, picture.LargeId)
			// Write media file description to a text file
			textFilePath := fmt.Sprintf("%s.txt", filePath)
			if _, err := os.Stat(textFilePath); errors.Is(err, os.ErrNotExist) {
				file, err := os.Create(textFilePath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to create file: %v\n", err)
					return
				}
				defer file.Close()

				fmt.Fprintln(file, picture.Label)
				fmt.Fprintln(file, picture.Description)
			}

			// Download media file itself
			fmt.Println(picture.ShortDate, picture.ImageLarge)

			extension := strings.Split(picture.Type, "/")[1]
			mediaFilePath := fmt.Sprintf("%s.%s", filePath, extension)
			if _, err := os.Stat(mediaFilePath); errors.Is(err, os.ErrNotExist) {
				downloadFile(mediaFilePath, picture.ImageLarge)

			}
			fyne.Do(func() {
				progressBar.SetValue(float64(i+1) / float64(totalPictures))
				progressBar.Refresh()
			})
		}
	}
}

type EBizPicture struct {
	Id                 int    `json:"id"`
	Label              string `json:"label"`
	Description        string `json:"description"`
	Date               string `json:"date"`
	LargeId            int    `json:"largeId"`
	AlbumId            int    `json:"albumId"`
	AlbumOrig          int    `json:"albumOrig"`
	CanEdit            bool   `json:"canEdit"`
	ShortDate          string `json:"shortDate"`
	ImgMedium          string `json:"imgMedium"`
	ImgMediumOrigin    string `json:"imgMediumOrigin"`
	ImgMediumSignedUrl string `json:"imgMediumSignedUrl"`
	ImgMediumId        int    `json:"imgMediumId"`
	ImgLarge           string `json:"imgLarge"`
	ImgLargeId         int    `json:"imgLargeId"`
	Image              string `json:"image"`
	ImageLarge         string `json:"imageLarge"`
	CanCommentOnPic    bool   `json:"canCommentOnPic"`
	UploadedBy         string `json:"uploadedBy"`
	UploadedById       string `json:"uploadedById"` // just... why the hell is this a string?
	IsAvailableTutor   bool   `json:"isAvailableTutor"`
	IsVideo            bool   `json:"isVideo"`
	Type               string `json:"type"`
	// pbo_context
	// pbo_routinetime
	// pbo_childcomments
	// pbo_obs
	// pbo_nextsteps
	// comments
}

type EBizPictures struct {
	Pictures []EBizPicture `json:"pictures"`
}

type EBizAuthenticate struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

type Child struct {
	Id   int
	Name string
}

func getChildPhotos(client http.Client, hostname string, child Child) []EBizPicture {
	var pictures []EBizPicture
	page := 1
	for {
		resp, _ := client.PostForm(
			fmt.Sprintf("%schildctrl/childgalleryloadmore", hostname),
			url.Values{
				"childId": {fmt.Sprint(child.Id)},
				"page":    {fmt.Sprint(page)},
			},
		)
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)

		var result EBizPictures
		if err := json.Unmarshal(respBody, &result); err != nil {
			fmt.Println("Can not unmarshal JSON")
			return nil
		}
		if len(result.Pictures) == 0 {
			break
		}
		pictures = append(pictures, result.Pictures...)
		page++
	}
	return pictures
}

// Borrowed from https://stackoverflow.com/a/33853856
func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func getChildrenList(client http.Client, hostname string) []Child {
	resp, err := client.Get(fmt.Sprintf("%seducators/home", hostname))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var children []Child
	// Find all child IDs and Names
	doc.Find(".educator-dash-child-name").Each(func(i int, s *goquery.Selection) {
		aref := s.Find("a")
		// For each item found, get the childName and url
		childName := strings.TrimSpace(aref.Text())
		childUrl, _ := aref.Attr("href")
		// childUrl in the format /child/######/info
		childIdRegex := `^/child/([0-9]+)/info$`
		re := regexp.MustCompile(childIdRegex)
		matches := re.FindStringSubmatch(childUrl)
		childId, _ := strconv.Atoi(matches[1])
		fmt.Printf("%d: %s\n", childId, childName)
		children = append(children, Child{Id: childId, Name: childName})
	})
	return children
}

// The authenticity token is bundled in the PLAY_SESSION cookie so this is used to get it
func getAuthenticityToken(hostname string, jar *cookiejar.Jar) (string, error) {
	u, _ := url.Parse(hostname)
	for _, cookie := range jar.Cookies(u) {
		if cookie.Name == "PLAY_SESSION" {
			authenticityTokenRegex := `^[0-9a-f]+-___AT=([0-9a-f]+)$`
			re := regexp.MustCompile(authenticityTokenRegex)
			matches := re.FindStringSubmatch(cookie.Value)
			return matches[1], nil
		}
	}
	// FIXME: this is obviously wrong
	return "", nil
}
