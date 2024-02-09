hljs.highlightAll()

function readingTime() {
    const article = document.getElementById("article")
    let post = article.innerHTML
    
    const WORDS_PER_MINUTE = 200;
    let result = {
        wordCount: 0,
        readingTime: 0
    };    //Matches words
    //See
    //https://regex101.com/r/q2Kqjg/6
    const regex=/\w+/g;
    result.wordCount = (post || '').match(regex).length;
    
    result.readingTime = Math.ceil(result.wordCount / WORDS_PER_MINUTE);

    const read = document.getElementById("read-time")

    read.innerHTML = `${result.readingTime} min read`
};