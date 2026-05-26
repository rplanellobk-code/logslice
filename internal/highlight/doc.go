// Package highlight provides ANSI-colour highlighting of regexp matches
// within log lines.
//
// Basic usage:
//
//	h, err := highlight.New(`ERROR`, highlight.WithColour(highlight.Red))
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(h.Apply(line))
//
// The zero value is not usable; always construct via New.
package highlight
