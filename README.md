# LiseBot
I wrote LiseBot for a friend of mine who had a large list (over 
1500 URLs) of Bookmarked Tweets which she wanted to save 
permanently on her 
computer. This by default is not something which Twitter provides, 
nor do they include Bookmarks in the data export tool, which is 
why I wrote this little script. 

## How it works
**Step 1:**
Download all your bookmarked Tweets as a CSV file using
[Dewey](https://getdewey.co) (or if you're doing something else 
you can just provide a list of Tweet URLs). 

**Step 2:**
Remove the commas and headers from the CSV file, so it's just a 
long list of Tweet URLs line by line.

**Step 3:**
Name the file `tweets.txt` in the same directory as LiseBot's 
source code, and 
create the folder `tweets` in said directory.

**Step 4:**
Compile and run Lisebot with `go mod download` and `go run main.go` and your download will begin!

**Results:**
Within your `tweets` directory, you'll have a text file with each 
Tweet's content, plus a separate file for each image of the Tweet.

**Extra tips:**
You can also set the `PROXY` environment variable to an HTTP(S) 
proxy URL to send the bot's traffic through.