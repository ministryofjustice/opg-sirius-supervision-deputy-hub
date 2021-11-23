import './main.scss';
import GOVUKFrontend from 'govuk-frontend/govuk/all.js';
import MojBannerAutoHide from './javascript/moj-banner-auto-hide';
import accessibleAutocomplete from 'accessible-autocomplete/dist/accessible-autocomplete.min.js';

GOVUKFrontend.initAll();

MojBannerAutoHide(document.querySelector('.app-main-class'));

accessibleAutocomplete.enhanceSelectElement({selectElement: document.querySelector("#select-ecm"), defaultValue: ""})