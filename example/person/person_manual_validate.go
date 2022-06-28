package person

import "fmt"

func (p *Person) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("empty name")
	}
	return nil
}
