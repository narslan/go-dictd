package main

import (
	_ "fmt"
	"github.com/narslan/dictdserver/database"
)

func main() {

	s := `
# dictd configuration file.
# whipped up by michael conrad tilstra <michael@gentoo.org>

# Informational message

global {
    site site.info
}

# who's allowed.  You might want to change this.
access {
  allow *
}

# Dictionaries are listed below.
# The initrc script scans /usr/lib/dict and adds all of the dictionaries
# it finds here.
#
# The initrc script will delete everything after the the last line and
# replace it with what it finds.  So add all of your things above.
#
# If this is a problem for people, contact me and
# we can work out a different method.
#
#deu-eng.dict.dz  deu-kur.dict.dz  eng-deu.dict.dz  eng-tur.dict.dz  kur-deu.dict.dz  kur-eng.dict.dz  tur-eng.dict.dz
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
}

#LASTLINE
access {  
    allow *
}

 database web1913 {
    data         "web1913.dict"
    index        "web1913.index"
    index_suffix "web1913.indexs"
    index_word   "web1913.indexw"

    default_strategy "word"

    access {
        allow *
    }
}
	`
	database.Parse(s)
	//fmt.Printf("%s", t)
	//	for toks := range l.Items {
	//		fmt.Printf("%v", toks)
	//	}
}
