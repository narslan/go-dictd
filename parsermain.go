package main

import (
	"fmt"
	"github.com/narslan/dictdserver/database"
)

func main() {

	s := `
#deu-eng.index    deu-kur.index    eng-deu.index    eng-tur.index    kur-deu.index    kur-eng.index    tur-eng.index

database deu-eng {
 data /usr/share/dictd/deu-eng.dict.dz
 index /usr/share/dictd/deu-eng.index
}
database deu-kur {
 data /usr/share/dictd/deu-kur.dict.dz
 index /usr/share/dictd/deu-kur.index
}
database eng-deu {
 data /usr/share/dictd/eng-deu.dict.dz
 index /usr/share/dictd/eng-deu.index
}
database eng-tur {
 data /usr/share/dictd/eng-tur.dict.dz
 index /usr/share/dictd/eng-tur.index
}
database kur-deu {
 data /usr/share/dictd/kur-deu.dict.dz
 index /usr/share/dictd/kur-deu.index
}
database kur-eng {
 data /usr/share/dictd/kur-eng.dict.dz
 index /usr/share/dictd/kur-eng.index
}
database tur-eng {
 data /usr/share/dictd/tur-eng.dict.dz
 index /usr/share/dictd/tur-eng.index
}
database wn {
 data /usr/share/dictd/wn.dict.dz
 index /usr/share/dictd/wn.index
}`

	r := database.ParseDatabases(s)
	for k, v := range r {
		fmt.Printf("%d %s\n", k, v)
	}
	//fmt.Printf("%T\n", r)
}
