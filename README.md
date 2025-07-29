# Scraper for Educabiz
Educabiz is a teacher-to-parent communication platform.
The teachers can send messages, reports, photos, and videos through there.
Overall, the idea is quite good, but the implementation lacks in some aspects.
When downloading pictures and videos from the platform, I noticed two things:
- The bulk downloader was not functioning properly.
- Downloading one picture at a time would take forever, and the downloaded picture or video was scaled down and compressed beyond what is reasonable.

After looking into the traffic of the site, I noticed that there is actually a reference to the original media files in the requests, it's just not used by the frontend though.
So, I made this scraper to ease the life of anyone who wants to keep their children's pictures in good quality and with minimal effort.
I could have made it a command-line tool if it were just for myself or even a bash script, but I took this opportunity to make something that is useful also to other parents who lack the expertise to do this and to learn how to create a simple cross-platform GUI in Go.

## Como utilizar (PT)
1. Faz download da última versão do programa para o teu sistema operativo [aqui](https://github.com/mamoit/educabiz-scraper/releases) (`.exe` para windows)
2. Põe o programa numa pasta vazia, ele irá fazer download dos conteúdos do educabiz para essa mesma pasta.
3. Corre o programa
4. Preenche o subdomínio com o nome da escola, que é em letras mínusculas e é a primeira parte do URL do educabiz da tua escola que usas para aceder ao site via web.
5. Preenche o username e a password.
6. Carrega `download` e espera que acabe. Dependendo da quantidade de conteúdos no site pode demorar um bom bocado. Durante este processo podes abrir a pasta num explorador de ficheiros para ires vendo o que o programa vai tirando do site.

## Usage (EN)
1. Download the program to an empty folder.
3. Run the program.
4. Fill in the subdomain of your school, it is usually the name in lower case, and it is the first part of the URL you use to access the website.
5. Fill in the user and password.
6. Press "download" and wait for it to finish.
