{{ define "layout" }}
<!doctype html>

<html>
<head> <title>{{ .Title }}</title>
<!--[if lt IE 7 ]> <body class="ie6"> <![endif]-->
<!--[if IE 7 ]>    <body class="ie7"> <![endif]-->
<!--[if IE 8 ]>    <body class="ie8"> <![endif]-->
<!--[if IE 9 ]>    <body class="ie9"> <![endif<]-->
<!--[if (gt IE 9)|!(IE)]><!-->  <!--<![endif]-->

<link rel="shortcut icon" href="/favicon.ico">
<link rel="apple-touch-icon" href="/apple-touch-icon.png">

<link rel="stylesheet" type="text/css" href="/css/main.css" />
 </head>
<body>
{{template "content" . }}
</body>
</html>
{{end}}
