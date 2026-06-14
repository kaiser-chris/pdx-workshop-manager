package steam

type ApiLanguage string

const (
	Arabic     ApiLanguage = "arabic"
	Bulgarian  ApiLanguage = "bulgarian"
	Schinese   ApiLanguage = "schinese"
	Tchinese   ApiLanguage = "tchinese"
	Czech      ApiLanguage = "czech"
	Danish     ApiLanguage = "danish"
	Dutch      ApiLanguage = "dutch"
	English    ApiLanguage = "english"
	Finnish    ApiLanguage = "finnish"
	French     ApiLanguage = "french"
	German     ApiLanguage = "german"
	Greek      ApiLanguage = "greek"
	Hungarian  ApiLanguage = "hungarian"
	Indonesian ApiLanguage = "indonesian"
	Italian    ApiLanguage = "italian"
	Japanese   ApiLanguage = "japanese"
	Korean     ApiLanguage = "korean"
	Malay      ApiLanguage = "malay"
	Norwegian  ApiLanguage = "norwegian"
	Polish     ApiLanguage = "polish"
	Portuguese ApiLanguage = "portuguese"
	Brazilian  ApiLanguage = "brazilian"
	Romanian   ApiLanguage = "romanian"
	Russian    ApiLanguage = "russian"
	Spanish    ApiLanguage = "spanish"
	Latam      ApiLanguage = "latam"
	Swedish    ApiLanguage = "swedish"
	Thai       ApiLanguage = "thai"
	Turkish    ApiLanguage = "turkish"
	Ukrainian  ApiLanguage = "ukrainian"
	Vietnamese ApiLanguage = "vietnamese"
)

func (language ApiLanguage) GetString() string {
	return string(language)
}

var ApiLanguages = map[ApiLanguage]string{
	Arabic:     "Arabic",
	Bulgarian:  "Bulgarian",
	Schinese:   "Chinese (Simplified)",
	Tchinese:   "Chinese (Traditional)",
	Czech:      "Czech",
	Danish:     "Danish",
	Dutch:      "Dutch",
	English:    "English",
	Finnish:    "Finnish",
	French:     "French",
	German:     "German",
	Greek:      "Greek",
	Hungarian:  "Hungarian",
	Indonesian: "Indonesian",
	Italian:    "Italian",
	Japanese:   "Japanese",
	Korean:     "Korean",
	Malay:      "Malay",
	Norwegian:  "Norwegian",
	Polish:     "Polish",
	Portuguese: "Portuguese",
	Brazilian:  "Portuguese-Brazil",
	Romanian:   "Romanian",
	Russian:    "Russian",
	Spanish:    "Spanish-Spain",
	Latam:      "Spanish-Latin America",
	Swedish:    "Swedish",
	Thai:       "Thai",
	Turkish:    "Turkish",
	Ukrainian:  "Ukrainian",
	Vietnamese: "Vietnamese",
}
