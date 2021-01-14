package object

import (
	"errors"
	"fmt"
	"github.com/roelofruis/spullen"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Form struct {
	Id        spullen.ObjectId
	TimeAdded string

	Name       string
	Quantity   string
	Categories string
	Tags       string
	Properties string
	Hidden     string
	Notes      string
	Marked     string

	Errors map[string]string

	object *spullen.Object
}

func EmptyForm() *Form {
	return &Form{Quantity: "1"}
}

func (f *Form) FillFromRequest(r *http.Request) {
	f.Name = r.PostFormValue("name")
	f.Quantity = r.PostFormValue("quantity")
	f.Categories = r.PostFormValue("categories")
	f.Tags = r.PostFormValue("tags")
	f.Properties = r.PostFormValue("properties")
	f.Hidden = r.PostFormValue("hidden")
	f.Notes = r.PostFormValue("notes")
	f.Marked = r.PostFormValue("marked")
}

func FormFromObject(o *spullen.Object) *Form {
	var propertyStrings []string = nil
	for _, p := range o.Properties {
		propertyStrings = append(propertyStrings, p.Key+"="+p.Value)
	}

	var hidden = ""
	if o.Hidden {
		hidden = "true"
	}

	var marked = "false"
	if o.Marked {
		marked = "true"
	}
	return &Form{
		Id:         o.Id,
		TimeAdded:  strconv.FormatInt(o.Added.Unix(), 10),
		Name:       o.Name,
		Quantity:   strconv.FormatInt(int64(o.Quantity), 10),
		Categories: strings.Join(o.Categories, ","),
		Tags:       strings.Join(o.Tags, ","),
		Properties: strings.Join(propertyStrings, ","),
		Hidden:     hidden,
		Notes:      o.Notes,
		Marked:     marked,
	}
}

func (f *Form) GetObject() (*spullen.Object, error) {
	if f.object == nil {
		return nil, errors.New("form has not been validated")
	}

	return f.object, nil
}

func (f *Form) Validate() bool {
	f.Errors = make(map[string]string)

	if len(f.Id) != 16 {
		f.Errors["Id"] = "Id moet bestaan uit 16 tekens"
	}

	if len(f.Name) == 0 {
		f.Errors["Name"] = "Voer een naam in"
	}

	t, err := strconv.ParseInt(f.TimeAdded, 10, 64)
	if err != nil {
		f.Errors["TimeAdded"] = "Geen geldige Unix tijdwaarde"
	}

	q, err := strconv.ParseInt(f.Quantity, 10, 32)
	if err != nil {
		f.Errors["Quantity"] = "Aantal moet een getal zijn"
	}

	if q < 1 {
		f.Errors["Quantity"] = "Aantal moet minstens 1 zijn"
	}

	var categories []string
	if f.Categories == "" {
		categories = []string{""}
	} else {
		for _, c := range strings.Split(f.Categories, ",") {
			if len(c) == 0 {
				f.Errors["Categories"] = "Categorieën mag geen lege waardes bevatten"
				break
			}
			categories = append(categories, normalize(c))
		}
	}

	var tags []string
	if f.Tags == "" {
		tags = []string{""}
	} else {
		for _, t := range strings.Split(f.Tags, ",") {
			if len(t) == 0 {
				f.Errors["Tags"] = "Tags mag geen lege waardes bevatten"
				break
			}
			tags = append(tags, normalize(t))
		}
	}

	var properties []*spullen.Property
	if f.Properties == "" {
		properties = []*spullen.Property{}
	} else {
		for _, p := range strings.Split(f.Properties, ",") {
			if len(p) == 0 {
				f.Errors["Properties"] = "Eigenschappen mag geen lege waardes bevatten"
				break
			}

			keyValue := strings.Split(p, "=")
			if len(keyValue) != 2 {
				f.Errors["Properties"] = fmt.Sprintf("Ongeldige eigenschap '%s'", p)
				break
			}
			properties = append(properties, &spullen.Property{
				Key:   normalize(keyValue[0]),
				Value: normalize(keyValue[1]),
			})
		}
	}

	hidden, err := strconv.ParseBool(f.Hidden)
	if err != nil {
		f.Errors["Hidden"] = "Verborgen moet een geldige booleaanse waarde zijn"
	}

	marked, err := strconv.ParseBool(f.Marked)
	if err != nil {
		f.Errors["Marked"] = "Gemarkeerd moet een geldige booleaanse waarde zijn"
	}

	isValid := len(f.Errors) == 0
	if isValid {
		f.object = &spullen.Object{
			Id:         f.Id,
			Added:      time.Unix(t, 0),
			Name:       normalize(f.Name),
			Quantity:   int(q),
			Categories: categories,
			Tags:       tags,
			Properties: properties,
			Hidden:     hidden,
			Notes:      f.Notes,
			Marked:     marked,
		}
	}

	return isValid
}

func normalize(s string) string {
	return strings.ToLower(strings.Trim(s, " "))
}
