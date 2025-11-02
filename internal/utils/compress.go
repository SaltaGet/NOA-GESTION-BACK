package utils

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"io"

	"github.com/DanielChachagua/GestionCar/pkg/models"
)

// func CompressToBase64(input string) (string, error) {
// 	var b bytes.Buffer
// 	w := zlib.NewWriter(&b)
// 	_, err := w.Write([]byte(input))
// 	if err != nil {
// 		return "", models.ErrorResponse(500, "Error al comprimir", err)
// 	}
// 	if err := w.Close(); err != nil {
// 		return "", models.ErrorResponse(500, "Error al finalizar la compresión", err)
// 	}

// 	encoded := base64.StdEncoding.EncodeToString(b.Bytes())
// 	return encoded, nil
// }

// func DecompressFromBase64(input string) (string, error) {
// 	data, err := base64.StdEncoding.DecodeString(input)
// 	if err != nil {
// 		return "", models.ErrorResponse(500, "Error al descomprimir", err)
// 	}

// 	b := bytes.NewReader(data)
// 	r, err := zlib.NewReader(b)
// 	if err != nil {
// 		return "", models.ErrorResponse(500, "Error al descomprimir", err)
// 	}
// 	defer r.Close()

// 	var out bytes.Buffer
// 	_, err = io.Copy(&out, r)
// 	if err != nil {
// 		return "", models.ErrorResponse(500, "Error al descomprimir", err)
// 	}

// 	return out.String(), nil
// }

func CompressToBase64Bytes(input []byte) (string, error) {
	var b bytes.Buffer

	w := zlib.NewWriter(&b)
	_, err := w.Write(input)
	if err != nil {
		return "", models.ErrorResponse(500, "Error al comprimir", err)
	}

	if err := w.Close(); err != nil {
		return "", models.ErrorResponse(500, "Error al finalizar la compresión", err)
	}

	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

func DecompressFromBase64(input string) (string, error) {
	compressedData, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", models.ErrorResponse(400, "Error al decodificar base64", err)
	}

	b := bytes.NewReader(compressedData)
	r, err := zlib.NewReader(b)
	if err != nil {
		return "", models.ErrorResponse(400, "Error al descomprimir datos", err)
	}
	defer r.Close()

	uncompressedData, err := io.ReadAll(r)
	if err != nil {
		return "", models.ErrorResponse(500, "Error al leer datos descomprimidos", err)
	}

	return string(uncompressedData), nil
}
