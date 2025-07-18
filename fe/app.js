async function loadPosts() {
  const res = await fetch("/posts");
  const posts = await res.json();

  // Group posts by source (if no source, use 'Feed')
  const grouped = {};
  posts.forEach((post) => {
    const source = post.source || "Feed";
    if (!grouped[source]) grouped[source] = [];
    grouped[source].push(post);
  });

  const feedList = document.getElementById("feed-list");
  feedList.innerHTML = "";
  let firstArticle = null;

  Object.keys(grouped).forEach((source) => {
    const sourceDiv = document.createElement("div");
    sourceDiv.className = "feed-source";
    sourceDiv.textContent = source;
    feedList.appendChild(sourceDiv);

    grouped[source].forEach((post, idx) => {
      const a = document.createElement("a");
      a.className = "article-title";
      a.textContent = post.title;
      a.href = "#";
      a.onclick = (e) => {
        e.preventDefault();
        showArticle(post, a);
      };
      feedList.appendChild(a);
      if (!firstArticle) {
        firstArticle = { post, a };
      }
    });
  });

  // Show first article by default
  if (firstArticle) {
    showArticle(firstArticle.post, firstArticle.a);
  }
}

function showArticle(post, linkElem) {
  document
    .querySelectorAll(".article-title.active")
    .forEach((el) => el.classList.remove("active"));
  if (linkElem) linkElem.classList.add("active");
  const detail = document.getElementById("article-detail");
  detail.innerHTML = `
    <h2>${post.title}</h2>
    <a href="${post.link}" target="_blank">Read original</a>
    <div style="margin-top:1em;">${post.description}</div>
  `;
}

window.onload = loadPosts;
