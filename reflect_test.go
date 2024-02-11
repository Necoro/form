package form

import (
	"html/template"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// fields is where like 99% of the real work gets done, so most of the
// testing effort should be focused here. It is also very easy to
// test - just plug in values and verify that you get the expected
// field slice back.
func Test_fields(t *testing.T) {
	type address struct {
		Street1 string
	}
	var nilAddress *address
	type addressWithTags struct {
		Street1 string `form:"name=street"`
	}

	tests := []struct {
		name string
		arg  interface{}
		want []field
	}{
		{
			name: "simple and empty",
			arg: struct {
				Name string
			}{},
			want: []field{
				{
					Name:        "Name",
					Label:       "Name",
					Placeholder: "Name",
					Type:        "text",
					Value:       "",
					ReadOnly:    false,
				},
			},
		}, {
			name: "simple with value",
			arg: struct {
				Name string
			}{"Michael Scott"},
			want: []field{
				{
					Name:        "Name",
					Label:       "Name",
					Placeholder: "Name",
					Type:        "text",
					Value:       "Michael Scott",
					ReadOnly:    false,
				},
			},
		}, {
			name: "simple with ignored",
			arg: struct {
				Name    string
				Ignored string `form:"-"`
			}{"", "secret info"},
			want: []field{
				{
					Name:        "Name",
					Label:       "Name",
					Placeholder: "Name",
					Type:        "text",
					Value:       "",
					ReadOnly:    false,
				},
			},
		}, {
			name: "pointer to struct w/ val",
			arg:  &address{},
			want: []field{
				{
					Name:        "Street1",
					Label:       "Street1",
					Placeholder: "Street1",
					Type:        "text",
					Value:       "",
					ReadOnly:    false,
				},
			},
		}, {
			name: "nil pointer with type",
			arg:  nilAddress,
			want: []field{
				{
					Name:        "Street1",
					Label:       "Street1",
					Placeholder: "Street1",
					Type:        "text",
					Value:       "",
					ReadOnly:    false,
				},
			},
		}, {
			name: "nested simple",
			arg: struct {
				Name    string
				Address struct {
					Street1 string
				}
			}{},
			want: []field{
				{
					Name:        "Name",
					Label:       "Name",
					Placeholder: "Name",
					Type:        "text",
					Value:       "",
					ReadOnly:    false,
				}, {
					Name:        "Address.Street1",
					Label:       "Street1",
					Placeholder: "Street1",
					Type:        "text",
					Value:       "",
					ReadOnly:    false,
				},
			},
		}, {
			name: "nested with values",
			arg: struct {
				Name    string
				Address address
			}{
				Name:    "Michael Scott",
				Address: address{"123 Test St"},
			},
			want: []field{
				{
					Name:        "Name",
					Label:       "Name",
					Placeholder: "Name",
					Type:        "text",
					Value:       "Michael Scott",
					ReadOnly:    false,
				}, {
					Name:        "Address.Street1",
					Label:       "Street1",
					Placeholder: "Street1",
					Type:        "text",
					Value:       "123 Test St",
					ReadOnly:    false,
				},
			},
		}, {
			name: "nested with tags",
			arg: struct {
				Name     string `form:"label=Full Name;id=name"`
				Password string `form:"type=password;footer=Something super secret!"`
				Address  addressWithTags
			}{
				Name:    "Michael Scott",
				Address: addressWithTags{"123 Test St"},
			},
			want: []field{
				{
					Name:        "Name",
					Label:       "Full Name",
					Placeholder: "Full Name",
					Type:        "text",
					Value:       "Michael Scott",
					ID:          "name",
					ReadOnly:    false,
				}, {
					Name:        "Password",
					Label:       "Password",
					Placeholder: "Password",
					Type:        "password",
					Value:       "",
					Footer:      template.HTML("Something super secret!"),
					ReadOnly:    false,
				}, {
					Name:        "street",
					Label:       "Street1",
					Placeholder: "Street1",
					Type:        "text",
					Value:       "123 Test St",
					ReadOnly:    false,
				},
			},
		}, {
			name: "nested with nil ptr",
			arg: struct {
				Name    string
				Address *address
			}{
				Name:    "Michael Scott",
				Address: nil,
			},
			want: []field{
				{
					Name:        "Name",
					Label:       "Name",
					Placeholder: "Name",
					Type:        "text",
					Value:       "Michael Scott",
				}, {
					Name:        "Address.Street1",
					Label:       "Street1",
					Placeholder: "Street1",
					Type:        "text",
					Value:       "",
				},
			},
		},
		{
			name: "custom css class",
			arg: struct {
				Name string `form:"class=custom-css-class"`
			}{
				Name: "Michael Scott",
			},
			want: []field{
				{
					Name:        "Name",
					Label:       "Name",
					Placeholder: "Name",
					Type:        "text",
					Value:       "Michael Scott",
					Class:       "custom-css-class",
					ReadOnly:    false,
				},
			},
		}, {
			name: "read only field",
			arg: struct {
				Name string `form:"readonly=true"`
			}{
				Name: "Michael Scott",
			},
			want: []field{
				{
					Name:        "Name",
					Label:       "Name",
					Placeholder: "Name",
					Type:        "text",
					Value:       "Michael Scott",
					ReadOnly:    true,
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := fields(tc.arg)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
