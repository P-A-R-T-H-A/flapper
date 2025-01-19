function init() {
    const toggleButton = document.getElementById('theme-toggle');
    const r = document.querySelector(':root');
    toggleButton.addEventListener('click', () => {
      let rs = getComputedStyle(r);
      if (rs.getPropertyValue('--primary-color') === '#fff') {
        r.style.setProperty('--primary-color', '#333'); 
      } else {
        r.style.setProperty('--primary-color', '#fff'); 
      }
    });
}

document.addEventListener("DOMContentLoaded", init);