package book

import (
	"io"
	"strings"
	"testing"

	"github.com/pirmd/verify"
)

func TestGetRawTextFromHTML(t *testing.T) {
	testCases := []struct {
		in   string
		want string
	}{
		{
			in:   "<p>Hello Gophers!</p>",
			want: "Hello Gophers!",
		},
		{
			in: `
			<div>
				<p>Hello Gophers!</p>
				<p>Golang is nice</p>
			</div>`,
			want: "Hello Gophers!\nGolang is nice",
		},
		{
			in:   "<p>Hello <span>Gophers</span>!</p>",
			want: "Hello Gophers!",
		},
		{
			in:   `Hello <b>Gophers</b>!`,
			want: `Hello Gophers!`,
		},
		{
			in:   `Hello <i>Gophers</i>!`,
			want: `Hello Gophers!`,
		},
		{
			in:   `Hello <del>Gophers</del>!`,
			want: `Hello Gophers!`,
		},
		{
			in:   `<h1>Hello <i>Gophers</i>!</h1>`,
			want: "Hello Gophers!",
		},
		{
			in:   `Hello<b> Gophers</b>!`,
			want: `Hello Gophers!`,
		},
		{
			in:   `Hello  <b> Gophers</b>!`,
			want: `Hello Gophers!`,
		},
		{
			in:   `  Hello  <b> Gophers</b>!`,
			want: `Hello Gophers!`,
		},
		{
			in:   `    <b>Hello Gophers</b>!`,
			want: `Hello Gophers!`,
		},
		{
			in:   `<a href="http://interesting.com/">Link</a>`,
			want: `Link`,
		},
		{
			in:   `<a onclick="alert(42)">Link</a>`,
			want: `Link`,
		},
		{
			in: `
            <ul>
                <li>todo</li>
                <li>really need todo</li>
            </ul>`,
			want: "todo\nreally need todo",
		},
		{
			in: `
            <ol>
                <li>first thing</li>
                <li>second thing</li>
            </ol>`,
			want: "first thing\nsecond thing",
		},
		{
			in: `
            <ul>
            <li>item1
                <ol>
                    <li>item1.1</li>
                    <li>item1.2</li>
                    <li>item1.3</li>
                    <li>item1.4</li>
                 </ol>
            </li>
            <li>item2
                 <ul>
                    <li>item2.1</li>
                    <li>item2.2</li>
                 </ul>
            </li>
            </ul>`,
			want: "item1\nitem1.1\nitem1.2\nitem1.3\nitem1.4\nitem2\nitem2.1\nitem2.2",
		},
		{
			in: `
            <table>
            <tr>
                <th>Col1</th>
                <th>Col2</th>
            </tr>
            <tr>
                <td>Col1.1</td>
                <td>Col2.1</td>
            </tr>
            <tr>
                <td>Col1.2</td>
                <td>Col2.2</td>
            </tr>
            </table>`,
			want: "Col1\nCol2\nCol1.1\nCol2.1\nCol1.2\nCol2.2",
		},
		{
			in: `
            <ul>
                <li>Hello Gophers!</li>
                <li>Golang is <b> so</b> nice</li>
            </ul>

            <script type='text/javascript'>
            really_useful_stuff();
            </script>`,
			want: "Hello Gophers!\nGolang is so nice",
		},
		{
			in: `
            <p>Hello Gophers!</p>
            <span>Golang is nice</span>`,
			want: "Hello Gophers!\nGolang is nice",
		},
		{
			in:   "<p>  Hello\n  <b> Gophers</b>!</p>",
			want: `Hello Gophers!`,
		},
		{
			in:   `<p class="UUID"><span class="isbn">ISBN 978-0-596-</span><span class="isbn">52068-7</span></p>`,
			want: `ISBN 978-0-596-52068-7`,
		},
		{
			in:   `<p class="p1">11 &ndash; X A12616 ISBN 978-2-07-012616-3&nbsp;13,90€</p>`,
			want: `11 – X A12616 ISBN 978-2-07-012616-3 13,90€`,
		},
	}

	for _, tc := range testCases {
		inR := strings.NewReader(tc.in)

		gotR, err := getRawTextFromHTML(inR)
		if err != nil {
			t.Errorf("Fail to extract text from '%v': %v", tc.in, err)
		}

		got, err := io.ReadAll(gotR)
		if err != nil {
			t.Errorf("Fail to extract text from '%v': %v", tc.in, err)
		}

		if failure := verify.Equal(string(got), tc.want); failure != nil {
			t.Errorf("Fail to extract text from '%v':\n%v", tc.in, failure)
		}
	}
}
