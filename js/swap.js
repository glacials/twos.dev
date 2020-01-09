document.addEventListener('DOMContentLoaded', function(event) {
  const bootstrapCss = document.getElementById('bootstrapCss')
  const txtCss = document.getElementById('txtCss')

  document.getElementById('toggleView').addEventListener('click', function(event) {
    bootstrapCss.disabled ^= true
    txtCss.disabled ^= true
    event.target.innerText = `view ${bootstrapCss.disabled ? 'rich' : 'text only'}`
  })
})
