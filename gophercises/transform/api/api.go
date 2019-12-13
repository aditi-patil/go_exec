package api

import (
	"errors"
	"fmt"
	"gophercises/transform/primitive"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
)

var newTempfile = tempfile
var ioCopy = io.Copy
var ioTempFile = ioutil.TempFile

// ModifyImage gives images with modified version having different modes and
// number of shapes.
func ModifyImage(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("./img/" + filepath.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer f.Close()
	ext := filepath.Ext(f.Name())[1:]
	modeStr := r.FormValue("mode")
	if modeStr == "" {
		renderModeChoices(w, r, f, ext)
		return
	}
	mode, err := strconv.Atoi(modeStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	nStr := r.FormValue("n")
	if nStr == "" {
		renderNumShapeChoices(w, r, f, ext, primitive.Mode(mode))
		return
	}
	numShapes, err := strconv.Atoi(nStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = numShapes
	http.Redirect(w, r, "/img/"+filepath.Base(f.Name()), http.StatusFound)
}

// UploadImage uploads image and redirect to the modify image path
func UploadImage(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	ext := filepath.Ext(header.Filename)[1:]
	onDisk, err := newTempfile("", ext)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	defer onDisk.Close()
	_, err = ioCopy(onDisk, file)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/modify/"+filepath.Base(onDisk.Name()), http.StatusFound)
}

func renderNumShapeChoices(w http.ResponseWriter, r *http.Request, rs io.ReadSeeker, ext string, mode primitive.Mode) {
	opts := []genOpts{
		{N: 10, M: mode},
		{N: 20, M: mode},
		{N: 30, M: mode},
		{N: 40, M: mode},
	}
	imgs, err := genImages(rs, ext, opts...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := `<html><body>
			{{range .}}
				<a href="/modify/{{.Name}}?mode={{.Mode}}&n={{.NumShapes}}">
					<img style="width: 20%;" src="/img/{{.Name}}">
				</a>
			{{end}}
			</body></html>`
	tpl := template.Must(template.New("").Parse(html))
	type dataStruct struct {
		Name      string
		Mode      primitive.Mode
		NumShapes int
	}
	var data []dataStruct
	for i, img := range imgs {
		data = append(data, dataStruct{
			Name:      filepath.Base(img),
			Mode:      opts[i].M,
			NumShapes: opts[i].N,
		})
	}
	err = tpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func renderModeChoices(w http.ResponseWriter, r *http.Request, rs io.ReadSeeker, ext string) {
	opts := []genOpts{
		{N: 10, M: primitive.ModeCircle},
		{N: 10, M: primitive.ModeBeziers},
		{N: 10, M: primitive.ModePolygon},
		{N: 10, M: primitive.ModeCombo},
	}
	imgs, err := genImages(rs, ext, opts...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := `<html><body>
			{{range .}}
				<a href="/modify/{{.Name}}?mode={{.Mode}}">
					<img style="width: 20%;" src="/img/{{.Name}}">
				</a>
			{{end}}
			</body></html>`
	tpl := template.Must(template.New("").Parse(html))
	type dataStruct struct {
		Name string
		Mode primitive.Mode
	}
	var data []dataStruct
	for i, img := range imgs {
		data = append(data, dataStruct{
			Name: filepath.Base(img),
			Mode: opts[i].M,
		})
	}
	err = tpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

type genOpts struct {
	N int
	M primitive.Mode
}

func genImages(rs io.ReadSeeker, ext string, opts ...genOpts) ([]string, error) {
	var ret []string
	for _, opt := range opts {
		rs.Seek(0, 0)
		f, err := genImageFile(rs, ext, opt.N, opt.M)
		if err != nil {
			return nil, err
		}
		ret = append(ret, f)
	}
	return ret, nil
}

func genImageFile(r io.Reader, ext string, numShapes int, mode primitive.Mode) (string, error) {
	out, err := primitive.NewTransform(r, ext, numShapes, primitive.WithMode(mode))
	if err != nil {
		return "", err
	}
	outFile, err := newTempfile("", ext)
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	io.Copy(outFile, out)
	return outFile.Name(), nil
}

func tempfile(prefix, ext string) (*os.File, error) {
	in, err := ioTempFile("./img/", prefix)
	if err != nil {
		return nil, errors.New("main: failed to create temporary file")
	}
	defer os.Remove(in.Name())
	return os.Create(fmt.Sprintf("%s.%s", in.Name(), ext))
}
