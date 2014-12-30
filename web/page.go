package web


type Page struct {
    Title string
    Content string 
}


func NewPage() *Page {
    p := new(Page)
    return p

}

func (p *Page) SetTitle(title string) (error) {
    p.Title = title
    return nil
}

func (p *Page) SetContent(content string) (error) {
    p.Content = content
    return nil
}
