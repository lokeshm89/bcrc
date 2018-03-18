var express =require('express');
var app = express();
var path = require('path');
app.use(express.static(path.join(__dirname,'public')));
app.get('/',function(req,res){
	res.sendFile(path.join(__dirname + '/index.html'));
});
app.get('/home',function(req,res){
	res.sendFile(path.join(__dirname + '/milk-input.html'));
});
app.get('/processing',function(req,res){
	res.sendFile(path.join(__dirname + '/processing.html'));
});
app.get('/transport',function(req,res){
	res.sendFile(path.join(__dirname + '/transport.html'));
});
app.listen(8080);
