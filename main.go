package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

func convertHeicToPng(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")

	if err != nil {
		http.Error(w, "Erro ao receber o arquivo", http.StatusBadRequest)
		return
	}

	defer file.Close()

	// Salvando arquivo temporario
	tmpFile, err := ioutil.TempFile("", "uploaded-*.heic")

	if err != nil {
		http.Error(w, "Erro ao criar arquivo temporario", http.StatusInternalServerError)
		return
	}

	defer os.Remove(tmpFile.Name())

	// Copiar o conteudo do arquivo recebido para arquivo temporario
	//_, err = ioutil.ReadAll(file)
	_, err = io.Copy(tmpFile, file)

	if err != nil {
		http.Error(w, "Erro ao salvar arquivo temporario", http.StatusInternalServerError)
		return
	}

	// garantir que os dados estejam escritos no disco
	tmpFile.Close()

	// Definir o caminho de saida para o arquivo png
	outputPath := tmpFile.Name() + ".png"

	// Usando ImageMagick para converter de heic para png
	cmd := exec.Command("C:\\Program Files\\ImageMagick-7.1.1-Q16-HDRI\\magick.exe", tmpFile.Name(), outputPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Erro ao executar comando: %v, Saída: %s\n", err, string(output))
		http.Error(w, "Erro ao converter imagem", http.StatusInternalServerError)
		return
	}

	// Ler arquivo convertido
	convertedImage, err := os.Open(outputPath)

	if err != nil {
		http.Error(w, "Erro ao abrir imagem convertida", http.StatusInternalServerError)
		return
	}

	defer convertedImage.Close()

	// cabecalho de resposta
	w.Header().Set("Content-Type", "image/png")

	/*
		// escrever imagem png na resposta
		_, err = ioutil.ReadAll(convertedImage)

		if err != nil {
			http.Error(w, "Erro ao enviar imagem", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	*/
	http.ServeFile(w, r, outputPath)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/convert", convertHeicToPng).Methods("POST")

	fmt.Println("Rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
