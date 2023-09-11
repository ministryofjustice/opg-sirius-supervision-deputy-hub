import GOVUKFrontend from "govuk-frontend/govuk/all.js";
import MojBannerAutoHide from "./javascript/moj-banner-auto-hide";
import accessibleAutocomplete from "accessible-autocomplete";
import "opg-sirius-header/sirius-header.js";

GOVUKFrontend.initAll();

MojBannerAutoHide(document.querySelector(".app-main-class"));

document.body.className = ((document.body.className) ? document.body.className + ' js-enabled' : 'js-enabled');

if (document.querySelector("#select-ecm")) {
    accessibleAutocomplete.enhanceSelectElement({
        selectElement: document.querySelector("#select-ecm"),
        defaultValue: "",
    });
}

if (document.querySelector("#select-existing-firm")) {
    accessibleAutocomplete.enhanceSelectElement({
        selectElement: document.querySelector("#select-existing-firm"),
        defaultValue: "",
    });
}

if (document.querySelector("#f-existing-firm")) {
    document.getElementById("f-existing-firm").onclick = function () {
        toggleAutocompleteInput();
    };
}

if (document.querySelector("#assignedto-other")) {
    document.getElementById("assignedto-other").onclick = function () {
        toggleAutocompleteInput();
    };
}

function toggleAutocompleteInput() {
    document
        .getElementById("autocomplete-input")
        .classList.toggle("hide");
}

if (document.querySelector("#f-back-button")) {
    document.getElementById("f-back-button").onclick = function (e) {
        e.preventDefault();
        history.go(parseInt(sessionStorage.getItem("backIndex")));
    }
}

if (document.querySelector("#f-button-disabled")) {
    document.getElementById("f-button-disabled").onclick = function (e) {
        e.preventDefault();
        document.getElementById("f-button-disabled-warning").classList.remove("hide");
    }
}

document.querySelectorAll(".min-date-today")
    .forEach(function(input) {
    input.setAttribute("min", new Date().toISOString().split('T')[0]);
});

function onHomePage() {
    const homePageUrlRegex = new RegExp('\\/\\d+\\/*$');
    return homePageUrlRegex.test(window.location.pathname);
}

function storeBackSessionVars(backIndex, href) {
    if (backIndex !== null && location.href === href) {
        sessionStorage.setItem("backIndex", (parseInt(backIndex) - 1).toString());
    }

    if (backIndex === null || href === null || location.href !== href || onHomePage()) {
        sessionStorage.setItem("backIndex", "-1");
        sessionStorage.setItem("href", location.href);
    }
}

storeBackSessionVars(sessionStorage.getItem("backIndex"), sessionStorage.getItem("href"));
