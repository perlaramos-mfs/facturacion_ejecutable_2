package ServerFile

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
)

func MoveFile(fileName string, pathInit string, pathFin string) {
	// Rename and Remove a file
	// Using Rename() function
	Original_Path := pathInit + "/" + fileName
	existDir := crearDirectorioSiNoExiste(pathFin)
	if existDir {
		New_Path := pathFin + "/" + fileName
		e := os.Rename(Original_Path, New_Path)
		if e != nil {
			log.Fatal(e)
		}
		fmt.Println(fileName, " movido al directorio ", pathFin)
	}
}

func crearDirectorioSiNoExiste(directorio string) bool {
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		err = os.Mkdir(directorio, 0755)
		if err != nil {
			// Aquí puedes manejar mejor el error, es un ejemplo
			fmt.Println("No se ha podido crear el directorio")
			return false
		}
	}
	return true
}

func ListarArchivo(directorio string) []fs.FileInfo {
	archivos, err := ioutil.ReadDir("/" + directorio)
	if err != nil {
		fmt.Println("No se ha podido listar la cantidad de archivos")
		return []fs.FileInfo{}
	}
	return archivos
}
