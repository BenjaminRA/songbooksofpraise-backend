package models

import (
	"fmt"

	"github.com/BenjaminRA/himnario-backend/db/sqlite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Parrafo struct {
	ID      int    `json:"id,omitempty"`
	Coro    bool   `json:"coro"`
	Parrafo string `json:"parrafo"`
	Acordes string `json:"acordes"`
}

type Tema struct {
	ID         int                `json:"id"`
	Tema       string             `json:"tema"`
	InsertedID primitive.ObjectID `json:"_id"`
	Himnos     []Himno
	SubTemas   []Tema
}

type Himno struct {
	ID       int    `json:"id,omitempty"`
	Titulo   string `json:"titulo"`
	Parrafos []Parrafo
	Temas    []Tema
}

func (n *Himno) GetHimnos() ([]Himno, error) {
	db := sqlite.GetDBConnection()
	q := `SELECT id, titulo FROM himnos where id <= 517`
	// Ejecutamos la query
	rows, err := db.Query(q)
	if err != nil {
		panic(err)
		return []Himno{}, err
	}

	defer rows.Close()
	himnos := []Himno{}
	for rows.Next() {
		rows.Scan(
			&n.ID,
			&n.Titulo,
		)
		himnos = append(himnos, *n)
	}
	return himnos, nil
}

func (n *Himno) GetCoros() ([]Himno, error) {
	db := sqlite.GetDBConnection()
	q := `SELECT id, titulo FROM himnos where id > 517`
	// Ejecutamos la query
	rows, err := db.Query(q)
	if err != nil {
		return []Himno{}, err
	}

	defer rows.Close()
	himnos := []Himno{}
	for rows.Next() {
		rows.Scan(
			&n.ID,
			&n.Titulo,
		)
		himnos = append(himnos, *n)
	}
	return himnos, nil
}

func (n *Parrafo) GetParrafos(himno_id int) ([]Parrafo, error) {
	db := sqlite.GetDBConnection()

	q := fmt.Sprintf("SELECT id, coro, parrafo, acordes FROM parrafos where himno_id = %v order by id", himno_id)

	rows, err := db.Query(q)
	if err != nil {
		return []Parrafo{}, err
	}

	defer rows.Close()

	parrafos := []Parrafo{}

	for rows.Next() {
		rows.Scan(
			&n.ID,
			&n.Coro,
			&n.Parrafo,
			&n.Acordes,
		)

		parrafos = append(parrafos, *n)
	}

	return parrafos, nil
}

func (n *Tema) GetTemas(himno_id int) ([]Tema, error) {
	db := sqlite.GetDBConnection()

	q := fmt.Sprintf("SELECT t.id, tema FROM temas t join tema_himnos th on t.id = th.tema_id where th.himno_id =%v", himno_id)

	rows, err := db.Query(q)

	if err != nil {
		return []Tema{}, err
	}

	temas := []Tema{}

	for rows.Next() {
		rows.Scan(
			&n.ID,
			&n.Tema,
		)
		n.SubTemas, _ = GetSubTemas(himno_id, n.ID)
		temas = append(temas, *n)
	}

	return temas, nil
}

func (n *Tema) GetAllTemas() ([]Tema, error) {
	db := sqlite.GetDBConnection()

	q := "SELECT id, tema from temas"

	rows, err := db.Query(q)

	if err != nil {
		return []Tema{}, err
	}

	temas := []Tema{}

	for rows.Next() {
		rows.Scan(
			&n.ID,
			&n.Tema,
		)
		temas = append(temas, *n)
	}

	return temas, nil
}

func (n *Tema) GetHimnos() ([]Himno, error) {
	db := sqlite.GetDBConnection()

	q := fmt.Sprintf("SELECT h.id, h.titulo from himnos h join tema_himnos th on h.id = th.himno_id where th.tema_id = %v", n.ID)

	rows, err := db.Query(q)

	if err != nil {
		return []Himno{}, err
	}

	himnos := []Himno{}
	himno := Himno{}
	for rows.Next() {
		rows.Scan(
			&himno.ID,
			&himno.Titulo,
		)
		himnos = append(himnos, himno)
	}

	return himnos, nil
}

func (n *Tema) GetSubTemaHimnos() ([]Himno, error) {
	db := sqlite.GetDBConnection()

	q := fmt.Sprintf("SELECT h.id, h.titulo from himnos h join sub_tema_himnos sth on h.id = sth.himno_id where sth.sub_tema_id = %v", n.ID)
	rows, err := db.Query(q)

	if err != nil {
		return []Himno{}, err
	}

	himnos := []Himno{}
	himno := Himno{}
	for rows.Next() {
		rows.Scan(
			&himno.ID,
			&himno.Titulo,
		)
		himnos = append(himnos, himno)
	}

	return himnos, nil
}

func (n *Tema) GetAllSubTemas(tema_id int) ([]Tema, error) {
	db := sqlite.GetDBConnection()

	q := fmt.Sprintf("SELECT id, sub_tema from sub_temas where tema_id = %v", tema_id)

	rows, err := db.Query(q)

	if err != nil {
		return []Tema{}, err
	}

	temas := []Tema{}
	tema := Tema{}
	for rows.Next() {
		rows.Scan(
			&tema.ID,
			&tema.Tema,
		)
		temas = append(temas, tema)
	}

	return temas, nil
}

func GetSubTemas(himno_id int, tema_id int) ([]Tema, error) {
	db := sqlite.GetDBConnection()

	q := fmt.Sprintf("SELECT st.id, st.sub_tema FROM sub_temas st join sub_tema_himnos sth on st.id = sth.sub_tema_id where sth.himno_id = %v and st.tema_id = %v", himno_id, tema_id)

	rows, err := db.Query(q)

	if err != nil {
		return []Tema{}, err
	}

	temas := []Tema{}

	n := Tema{}
	for rows.Next() {
		rows.Scan(
			&n.ID,
			&n.Tema,
		)
		temas = append(temas, n)
	}

	return temas, nil
}
