var express =require('express');
var app = express();
var path = require('path');
app.use(express.static(path.join(__dirname,'public')));
app.get('/home',function(req,res){
	res.sendFile(path.join(__dirname + '/index.html'));
});
app.get('/:batchID',function(req,res){
	res.sendFile(path.join(__dirname + '/history.html'));
});
app.listen(8085);
