// cO(alg) Interactive Enhancements & UX Utilities

document.addEventListener("DOMContentLoaded", function () {
  initCodeCopyButtons();
  initSidebarActiveState();
});

// HTMX event hook to re-initialize after dynamic content swaps
document.body.addEventListener("htmx:afterSwap", function (evt) {
  initCodeCopyButtons();
  initSidebarActiveState();
});

/**
 * Automatically injects a "Copiază codul" (Copy code) button on all code snippets
 */
function initCodeCopyButtons() {
  const codeBlocks = document.querySelectorAll(
    '.lesson-prose pre, div[style*="background: #0d1117"] pre, div[style*="background:#0d1117"] pre'
  );

  codeBlocks.forEach(function (pre) {
    const parent = pre.parentElement;
    if (parent.querySelector(".code-copy-btn")) return; // already added

    // Ensure parent has position relative
    parent.style.position = "relative";

    const copyBtn = document.createElement("button");
    copyBtn.className = "code-copy-btn";
    copyBtn.type = "button";
    copyBtn.innerHTML = `
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
        <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
      </svg>
      <span>Copiază</span>
    `;

    copyBtn.addEventListener("click", function () {
      const codeText = pre.innerText || pre.textContent;
      navigator.clipboard.writeText(codeText).then(function () {
        copyBtn.classList.add("copied");
        copyBtn.innerHTML = `
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
          <span>Copiat!</span>
        `;
        setTimeout(function () {
          copyBtn.classList.remove("copied");
          copyBtn.innerHTML = `
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
            </svg>
            <span>Copiază</span>
          `;
        }, 2200);
      }).catch(function (err) {
        console.error("Nu s-a putut copia codul: ", err);
      });
    });

    parent.appendChild(copyBtn);
  });
}

/**
 * Filters lessons in the sidebar dynamically as the user types
 */
function filterSidebarLessons(query) {
  const q = (query || "").toLowerCase().trim();
  const chapterGroups = document.querySelectorAll(".menu-chapter-group");

  chapterGroups.forEach(function (group) {
    let hasVisibleLesson = false;
    const lessons = group.querySelectorAll(".less");

    lessons.forEach(function (less) {
      const text = (less.textContent || "").toLowerCase();
      if (text.includes(q)) {
        less.style.display = "flex";
        hasVisibleLesson = true;
      } else {
        less.style.display = "none";
      }
    });

    // Also check chapter title itself
    const chapterTitle = (group.querySelector(".capitol")?.textContent || "").toLowerCase();
    if (chapterTitle.includes(q)) {
      hasVisibleLesson = true;
      lessons.forEach(l => l.style.display = "flex");
    }

    if (hasVisibleLesson || q === "") {
      group.style.display = "block";
    } else {
      group.style.display = "none";
    }
  });
}

/**
 * Highlights active lesson matching the URL query param
 */
function initSidebarActiveState() {
  const urlParams = new URLSearchParams(window.location.search);
  const activeId = urlParams.get("id");
  if (!activeId) return;

  const links = document.querySelectorAll(".menu .less a");
  links.forEach(function (link) {
    const href = link.getAttribute("href") || "";
    if (href.includes("id=" + activeId)) {
      link.parentElement.classList.add("active-lesson");
    } else {
      link.parentElement.classList.remove("active-lesson");
    }
  });
}

/**
 * Backward compatibility for old script calls if any
 */
function show_menu() {
  const menu = document.getElementById("menu");
  if (menu) {
    menu.classList.toggle("open");
  }
}

function slider_pos() {
  // deprecated, no-op for backward compatibility
}
