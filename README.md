# Gator<br>
Gator is an rss fetcher that brings you posts from your favorite feeds on the warm embrace of your terminal.
![gatorlogo](https://github.com/user-attachments/assets/e6ca413d-d33f-41bd-bf3b-7f83a9bbc835)
## Requirements:
You will need Go 1.2 or later, You can install Golang with webi or [Follow the official install instructions](https://go.dev/doc/install).
### Windows:
```
curl.exe https://webi.ms/golang | powershell
```
### Linux:
```
curl -sS https://webi.sh/golang | sh; \
source ~/.config/envman/PATH.env
```
Now that you have go the next requirement would be postgresql to store your feeds and posts,they have a nice install guide on their page:
### [Download postgresql](https://www.postgresql.org/download/linux/ubuntu/)
And for the last requirement there's [goose](https://github.com/pressly/goose) this shouldn't be a requirement but im kinda lazy, maybe next time i will use it to build the database itself, the thing is you download it with go:
```
go install https://github.com/pressly/goose
```
and i will explain how to use it in the setup process.
## Installing:
### To install gator you have to use the go package manager:
```
go install github.com/Cacutss/gator
```
Now you can use gator on your terminal and it will tell you the current version, but to truly use the program in it's intended purpose you have to set up the database
## How to set up database:
if you are on linux you may have to set up a password this may or not be required, the command would be like this:
```
sudo passwd postgres
```
you can set whatever password you want
then you have to do a quick little 
```
sudo -u postgres psql
```
inside postgres you do this, you can use whatever password you want but you have to remember it
```
CREATE DATABASE gator;
\c gator
ALTER USER postgres PASSWORD 'yourpassword';
```
then you can just type exit to get out of postgres
Now that the databse is all set up
You need goose to fill the database with all the tables that will hold feeds,posts and users.
To use goose first check if it's installed by using
```
goose -version
```
then navigate to the sql/schema directorie and you're gonna have to do create your own "key" to access the database, this will be used 2 times so for convenience i would keep it in the clipboard for a little bit. The command will be the following and you're gonna have to replace the username and password.
```
goose postgres "postgres://username:password@localhost:5432/gator" up
```
and then you have to use that same string(with your username:password) and use it with
```
gator setdb <your db url>
```
Once done you are ready to use gator!
I'll leave some commands and descriptions here so you can start using it:
## Commands:
* "gator register \<username\>" registers your user and logs you in
* "gator login \<username\>" logs in as that username
* "gator addfeed \<feedname\> \<url\>" adds a feed to the database and makes you follow it (following makes you able to get that feed post's)
* "gator feeds" shows you a list of the feeds and who added them.
* "gator follow \<url\>" makes the current logged in user follow a feed.
* "gator following" shows you a list of the feeds you are currently following.
* "gator unfollow \<url\>" Unfollows that feed.
* "gator users" shows you a list of all the users and marks whoever is logged in at the moment.
* "gator agg \<time\>" first of all time can be as an example: ("10s") == 10 seconds, or ("3m10s") = 3 minutes 10 seconds, yada yada. The point is it searchs for posts within the feeds in \<time\> interval.
* "gator browse \<number\>" shows you up to <number> posts.
* "gator reset" SHOULD NOT BE USED IN NORMAL CIRCUNSTANCES, removes all users.
Hope the quality is at least acceptable 
