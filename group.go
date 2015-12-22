package i18n

import (
	"fmt"

	"golang.org/x/text/language"
)

// Group is a translation subgroup
type Group struct {
	i18n  *I18n
	group string
}

func (group *Group) key(k string) string {
	return fmt.Sprintf("%s.%s", group.group, k)
}

// Group gets a translation group
func (group *Group) Group(key string) *Group {
	return &Group{
		i18n:  group.i18n,
		group: group.key(key),
	}
}

// GenerateHelper generates a method that allways gets tags in a certain language
// This is usefull for passing to the template engine
func (group *Group) GenerateHelper(tag language.Tag) T {
	return T(func(key string) string {
		return group.T(tag, group.key(key))
	})
}

// T is a helper method to get translation by lang string or language tag
func (group *Group) T(lang interface{}, key string) string {
	return group.i18n.T(lang, group.key(key))
}

// GetWithLangString parses the lang string before lookip up the translation
func (group *Group) GetWithLangString(lang string, key string) (*Translation, error) {
	return group.i18n.GetWithLangString(lang, group.key(key))
}

// Get translation
func (group *Group) Get(lang language.Tag, key string) *Translation {
	return group.i18n.Get(lang, group.key(key))
}

// Add translation
func (group *Group) Add(translation *Translation) error {
	translation.Key = group.key(translation.Key)
	return group.i18n.Add(translation)
}

// Delete translation
func (group *Group) Delete(translation *Translation) error {
	return group.i18n.Delete(translation)
}
