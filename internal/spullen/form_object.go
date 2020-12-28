package spullen

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ObjectForm struct {
	Id        string
	TimeAdded string

	Name       string
	Quantity   string
	Categories string
	Tags       string
	Properties string
	Hidden     string
	Notes      string

	Errors map[string]string

	object *Object
}

func EmptyForm() *ObjectForm {
	return &ObjectForm{Quantity: "1"}
}

func FormFromObject(o *Object) *ObjectForm {
	var propertyStrings []string = nil
	for _, p := range o.Properties {
		propertyStrings = append(propertyStrings, p.Key+"="+p.Value)
	}

	var hidden = ""
	if o.Hidden {
		hidden = "true"
	}
	return &ObjectForm{
		Id:         o.Id,
		TimeAdded:  strconv.FormatInt(o.Added.Unix(), 10),
		Name:       o.Name,
		Quantity:   strconv.FormatInt(int64(o.Quantity), 10),
		Categories: strings.Join(o.Categories, ","),
		Tags:       strings.Join(o.Tags, ","),
		Properties: strings.Join(propertyStrings, ","),
		Hidden:     hidden,
		Notes:      o.Notes,
	}
}

func (f *ObjectForm) GetObject() (*Object, error) {
	if f.object == nil {
		return nil, errors.New("form has not been validated")
	}

	return f.object, nil
}

func (f *ObjectForm) Validate() bool {
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

	var categories []string
	if f.Categories == "" {
		categories = []string{""}
	} else {
		for _, c := range strings.Split(f.Categories, ",") {
			if len(c) == 0 {
				f.Errors["Categories"] = "Categorieën mag geen lege waardes bevatten"
				break
			}
			categories = append(categories, strings.ToLower(c))
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
			tags = append(tags, strings.ToLower(t))
		}
	}

	var properties []*Property
	if f.Properties == "" {
		properties = []*Property{}
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
			properties = append(properties, &Property{
				strings.ToLower(keyValue[0]),
				strings.ToLower(keyValue[1]),
			})
		}
	}

	hidden, err := strconv.ParseBool(f.Hidden)
	if err != nil {
		f.Errors["Hidden"] = "Verborgen moet een geldige booleaanse waarde zijn"
	}

	isValid := len(f.Errors) == 0
	if isValid {
		f.object = &Object{
			Id:         f.Id,
			Added:      time.Unix(t, 0),
			Name:       strings.ToLower(f.Name),
			Quantity:   int(q),
			Categories: categories,
			Tags:       tags,
			Properties: properties,
			Hidden:     hidden,
			Notes:      f.Notes,
		}
	}

	return isValid
}
