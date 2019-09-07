module github.com/raddare

go 1.12

require (
	gitlab.com/mikrowezel/config v0.0.0
	gitlab.com/mikrowezel/log v0.0.0
	gitlab.com/mikrowezel/service v0.0.0
)

replace gitlab.com/mikrowezel/log => ../log

replace gitlab.com/mikrowezel/config => ../config

replace gitlab.com/mikrowezel/service => ../service
