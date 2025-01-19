const toggleButton = document.getElementById("order-toggle");
const r = document.querySelector(":root");
let propertyTypeBoolean = true;

toggleButton.addEventListener("click", () => {
  console.log("Clicked with propertyTypeBoolean =", propertyTypeBoolean);

  if (propertyTypeBoolean) {
    r.style.setProperty("--nav-logo-order", "2");
    r.style.setProperty("--nav-invisible-order", "1");
    propertyTypeBoolean = false;
  } else {
    r.style.setProperty("--nav-logo-order", "1");
    r.style.setProperty("--nav-invisible-order", "2");
    propertyTypeBoolean = true;
  }
});
