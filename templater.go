package gotemplater

//Templater manage
type templater struct {
	namespaces map[string]string //namespace to path
}

func (tpt *templater) Debug(){

	fmt.Println("debug")
}
var Templater *templater = nil

//InitTemplater should call once at begining of app to init templater
func InitTemplater() *templater {
	if Templater == nil {
		Templater = &templater{
			namespaces: map[string]string{}
		}
	}
	return Templater
}
