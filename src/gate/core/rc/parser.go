package rc

import (
	"io"
)

type Section struct {
	Resources map[string]string
}

type File struct {
	Anonymous *Section
	Sections map[string]*Section
}

func (self *File) readSection(content *FileContent) (result *Section, err error) {
	//fmt.Printf("<readSection %s>\n", content.Debug())
	result = &Section{Resources: make(map[string]string)}
	var k rune
	for done := !content.IsValid(); !done; {
		k, err = content.Current()
		if err != nil {
			return
		}
		if k == '[' {
			done = true
		} else {
			var (
				key string
				value string
			)
			_, err = content.SkipBlanks()
			if err != nil {
				return
			}
			if content.IsValid() {
				key, err = content.SkipWord()
				if err != nil {
					return
				}
				_, err = content.SkipBlanks()
				if err != nil {
					return
				}
				_, err = content.SkipSymbol("=")
				if err != nil {
					return
				}
				_, err = content.SkipBlanks()
				if err != nil {
					return
				}
				value, err = content.SkipToEndOfLine()
				if err != nil {
					return
				}
				_, err = content.SkipBlanks()
				if err != nil {
					return
				}
				//fmt.Printf("%s = %s\n", key, value)
				result.Resources[key] = value
			}
			done = !content.IsValid()
		}
	}
	//fmt.Printf("</readSection %s>\n", content.Debug())
	return
}

func (self *File) readAnonymousSection(content *FileContent) (err error) {
	self.Anonymous, err = self.readSection(content)
	return
}

func (self *File) readNamedSection(content *FileContent) (err error) {
	_, err = content.SkipBlanks()
	if err != nil {
		return
	}
	if content.IsValid() {
		var (
			sectionName string
			section *Section
		)
		_, err = content.SkipSymbol("[")
		if err != nil {
			return
		}
		sectionName, err = content.SkipWord()
		if err != nil {
			return
		}
		_, err = content.SkipSymbol("]")
		if err != nil {
			return
		}
		_, err = content.SkipToEndOfLine()
		if err != nil {
			return
		}
		_, err = content.SkipBlanks()
		if err != nil {
			return
		}
		section, err = self.readSection(content)
		if err != nil {
			return
		}
		self.Sections[sectionName] = section
	}
	return
}

func Read(in io.Reader) (result *File, err error) {
	result = &File{Sections: make(map[string]*Section)}
	content := readFile(in)
	if content.IsValid() {
		err = result.readAnonymousSection(content)
	}
	for content.IsValid() && err == nil {
		err = result.readNamedSection(content)
	}
	return
}
