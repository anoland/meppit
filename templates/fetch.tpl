{{define "fetch"}}
{{template "header" .}}
      <div class="container">
	<ul>
	{{range  .items}}
	<li>{{.}}</li>
	{{end}}
	</ul>
</div>
{{template "footer" .}}
{{end}}
