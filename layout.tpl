{{ define "layout" }}
<!DOCTYPE html>
<html lang="en">
<head> 
  <meta charset="utf-8">
  <title>{{ .Title }}</title>
  <meta name="description" content="">
  <meta name="author" content="">
  <meta name="viewport" content="width=device-width, initial-scale=1">
<!--[if lt IE 7 ]> <body class="ie6"> <![endif]-->
<!--[if IE 7 ]>    <body class="ie7"> <![endif]-->
<!--[if IE 8 ]>    <body class="ie8"> <![endif]-->
<!--[if IE 9 ]>    <body class="ie9"> <![endif<]-->
<!--[if (gt IE 9)|!(IE)]><!-->  <!--<![endif]-->

  <link rel="shortcut icon" type="image/png" href="/images/favicon.png">
  <link rel="apple-touch-icon" href="/apple-touch-icon.png">
  <link href="//fonts.googleapis.com/css?family=Raleway:400,300,600" rel="stylesheet" type="text/css">

  <link rel="stylesheet" type="text/css" href="/css/normalize.css">
  <link rel="stylesheet" type="text/css" href="/css/skeleton.css">
  <link rel="stylesheet" type="text/css" href="/css/main.css" />
</head>
<body>
  <div class="container ">
     
    <div id="nav" class="row">
        <span class="five columns">
            <ul>
                <li><a href="/">Home</a></li>
                <li><a href="/#">About</a></li>
                <li><a href="/#">Contact</a></li>
            </ul>
        </span>
        <span class="one columns u-pull-right">
            <ul>
              <li><a href="/login">Login</a></li>
            </ul>
        </span>
    </div>
    
    <div class="row u-full-width">
        <div><span style="text-align:center"> header section</span></div>
    </div>
    <div class="row">
      <div class="one-third column" >
        <h4>Rolling div</h4>
        <p>{{template "content" . }}</p>
      </div>
      <div class="two-thirds column" >
        <h4>map div</h4>

        <p>{{template "content" . }}</p>
      </div>
    </div>
  </div>
</body>
</html>
{{end}}
