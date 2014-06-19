[[define "mgo"]]
[[template "head"]]
<div>[[.Count]] results.</div>
<table>
<tr><td>&nbsp</td><td>Hostname</td><td>Notify</td><td>Method</td><tr>
[[range .Results]]
<tr><td><a href="/user/edit?id=[[.Id|StringId]]">[edit]</a></td><td>[[.Name]] </td><td>[[.Phone|PrettyPrint]]</td><td>[[.Timestamp|SaneTime]]</td></tr>
[[end]]
<table>
[[template "foot"]]
[[ end ]]