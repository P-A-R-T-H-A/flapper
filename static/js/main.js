function init() {
    console.log("Hello Worldf");

    const toggleButton = document.getElementById('theme-toggle');
    const root = document.documentElement;
    console.log("toggleButton =", toggleButton);
    toggleButton.addEventListener('click', () => {
        console.log("Clicked");
        console.log("Here with", root.style.getPropertyValue('--primary-color'));
      if (root.style.getPropertyValue('--primary-color') === '#fff') {
        console.log("Here", root.style.getPropertyValue('--primary-color'));
        root.style.setProperty('--primary-color', '#333'); 
      } else {
        console.log("There");
        root.style.setProperty('--primary-color', '#fff'); 
      }
    });

}

document.addEventListener("DOMContentLoaded", init);