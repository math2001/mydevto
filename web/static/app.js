window.addEventListener("load", e => {

  document.querySelector("#login").addEventListener("click", e => {
    e.preventDefault()
    const win = window.open(e.target.href, 'Login with github', 'width=500,height=500,centerscreen')
    if (window.focus) { win.focus() }
  })

})
