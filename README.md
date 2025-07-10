# Scraper for Educabiz

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

## Compiling from linux
### Dependencies
Install fyne's dependencies https://docs.fyne.io/started/

`make` if you want to use the `Makefile`.

### Cross compiling to Windows
``` bash
make windows
```

### Cross compiling to Mac OSX
I'm not doing this.
Cross compiling to Mac OSX seems to be way too much of a hassel and seems to require an apple developer account.
If someone wants to contribute a solution to enable the generation of binaries for Mac OSX please leave me a PR.
