<html>
<center>
<head>
       <title>Upload CSV, Download TXT results</title>
</head>
<body>
<h1>Upload a CSV file, exported from generations</h1>
<h2>You will get a TXT file back, review that file for entry<h2>
<form enctype="multipart/form-data" action="/upload" method="post">
    <h3>Choose a csv to upload</h3>
    <input type="file" name="uploadfile" />
    <input type="hidden" name="token" value="{{.}}"/>
    <input type="submit" value="upload" />
</form>
</body>
</center>
</html>
