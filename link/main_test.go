package link

import (
	"fmt"
	"os"
	"testing"
)

func TestParseLinks(t *testing.T) {
	cases := []struct {
		file string
		want []Link
	}{
		{
			"ex1.html", []Link{
				{
					Href: "/other-page",
					Text: "A link to another page",
				},
			},
		},
		{
			"ex2.html", []Link{
				{
					Href: "https://www.twitter.com/joncalhoun",
					Text: "Check me out on twitter",
				},
				{
					Href: "https://github.com/gophercises",
					Text: "Gophercises is on Github!",
				},
			},
		},
		{
			"ex3.html", []Link{
				{
					Href: "#",
					Text: "Login",
				},
				{
					Href: "https://twitter.com/marcusolsson",
					Text: "@marcusolsson",
				},
				{
					Href: "/lost",
					Text: "Lost? Need help?",
				},
			},
		},
		{
			"ex4.html", []Link{
				{
					Href: "/dog-cat",
					Text: "dog cat",
				},
			},
		},
		{
			"ex5.html", []Link{
				{
					Href: "",
					Text: "Look! No href",
				},
			},
		},
		{
			"ex6.html", []Link{
				{
					Href: "",
					Text: "",
				},
				{
					Href: "",
					Text: "",
				},
				{
					Href: "",
					Text: "",
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("test %s", c.file), func(t *testing.T) {
			file, err := os.Open(c.file)
			got, err := ParseLinks(file)
			if err != nil {
				t.Fatal("Could not parse HTML links", err)
			}

			assertEquals(t, got, c.want)
		})
	}
}

func assertEquals(t testing.TB, got []Link, want []Link) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatal("Link slices are of different lengths")
	}
	for i, v := range got {
		if v != want[i] {
			t.Errorf("got[%d] = %v != want[%d] = %v", i, v, i, want[i])
		}
	}
}
