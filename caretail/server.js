var express =require('express');
var app = express();
var path = require('path');
app.use(express.static(path.join(__dirname,'public')));
app.get('/',function(req,res){
	res.sendFile(path.join(__dirname + '/quality-check.html'));
});
app.get('/home2',function(req,res){
	res.sendFile(path.join(__dirname + '/quality-check.html'));
});
app.listen(8082);
