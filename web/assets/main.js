import "./main.scss";
import GOVUKFrontend from "govuk-frontend/govuk/all.js";
import MojBannerAutoHide from "./javascript/moj-banner-auto-hide";
import accessibleAutocomplete from "accessible-autocomplete";

GOVUKFrontend.initAll();

MojBannerAutoHide(document.querySelector(".app-main-class"));

if (document.querySelector("#select-ecm")) {
    accessibleAutocomplete.enhanceSelectElement({
        selectElement: document.querySelector("#select-ecm"),
        defaultValue: "",
    });
}

export function downloadClientList(deputyId) {
    const baseUrl = document.querySelector('[name=api-base-uri]').getAttribute('content')
    console.log("baseUrl");
    console.log(baseUrl);
    console.log(baseUrl);
    console.log("deputyId");
    console.log(deputyId);
    fetch(`${baseUrl}/api/v1/deputies/${deputyId}/clients-list`, {
        method: "GET",
        credentials: 'include',
        headers: {
            "Content-type": "application/csv",
            "OPG-Bypass-Membrane": 1,
        }
    })
    .then((response) => {
        return response.json();
    })
}