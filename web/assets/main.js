import GOVUKFrontend from "govuk-frontend/govuk/all.js";
import MojBannerAutoHide from "./javascript/moj-banner-auto-hide";
import accessibleAutocomplete from "accessible-autocomplete";

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
        toggleChangeFirmAutoCompleteHideClass();
    };
}

function toggleChangeFirmAutoCompleteHideClass() {
    document
        .getElementById("change-firm-autocomplete")
        .classList.toggle("hide");
}

if (document.querySelector("#f-back-button")) {
    document.getElementById("f-back-button").onclick = function (e) {
        e.preventDefault();
        history.back();
    }
}

if (document.querySelector("#f-button-disabled")) {
    document.getElementById("f-button-disabled").onclick = function (e) {
        e.preventDefault();
        document.getElementById("f-button-disabled-warning").classList.remove("hide");
    }
}

function sorting(e) {
        const dataUrl = e.target.dataset.sort
        window.location.href = dataUrl;
}

document.addEventListener('click', (e) => {
    if (e.target) {
        const moduleName = e.target.dataset.module;
        if (e.target.matches('[data-module="callHref"]')) {
            switch (moduleName) {
                case "callHref":
                    sorting(e)
                    break;
            }
        }
    }
});
