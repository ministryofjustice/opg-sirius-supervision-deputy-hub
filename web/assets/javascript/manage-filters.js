export default class ManageFilters {
    constructor(element) {
        this.toggleFilterButton = element.querySelectorAll(".js-container-button");
        this.applyFiltersButton = element.querySelector('[data-module="apply-filters"]');
        this.clearFilters = element.querySelector('[data-module="clear-filters"]');
        this.filters = element.querySelectorAll('[data-module="filter"]');
        this.filterComponents = element.querySelectorAll(".moj-filter__options");

        this.setupEventListeners();
        this.setupFilters();
    };

    setupEventListeners() {
        this.toggleFilterButton.forEach((element) => {
            this.toggleFilter = this.toggleFilter.bind(this);
            element.addEventListener("click", this.toggleFilter);
        });

        this.applyFilters = this.applyFilters.bind(this);
        this.applyFiltersButton.addEventListener("click", this.applyFilters);

        this.filterComponents.forEach((element) => {
            this.toggleFilterVisibility = this.toggleFilterVisibility.bind(this);
            element
                .querySelectorAll(".filter-toggle-button")[0]
                .addEventListener("click", this.toggleFilterVisibility);
        });
    }

    setupFilters() {
        this.filterComponents.forEach((element) => {
            const filterName = element.dataset.filterName;
            let isOpen = this.getFilterStatusFromLocalStorage(filterName);

            this.setFilterVisibility(element, isOpen);
        });
    }

    setFilterStatusToLocalStorage(filterName, isOpen) {
        window.sessionStorage.setItem(
            filterName,
            JSON.stringify({ value: isOpen })
        );
    }

    getFilterStatusFromLocalStorage(filterName) {
        let sessionStorageValue = JSON.parse(
            window.sessionStorage.getItem(filterName)
        );
        if (!sessionStorageValue) {
            sessionStorageValue = { value: false };
            this.setFilterStatusToLocalStorage(
                filterName,
                sessionStorageValue.value
            );
        }

        return sessionStorageValue.value;
    }


    toggleFilterVisibility(e) {
        const filterElement = e.target.parentElement.parentElement.parentElement;
        const filterName = filterElement.dataset.filterName;

        let isOpen = this.getFilterStatusFromLocalStorage(filterName);
        isOpen = isOpen !== true;

        this.setFilterVisibility(filterElement, isOpen);

        this.setFilterStatusToLocalStorage(filterName, isOpen);
    }

    setFilterVisibility(element, isOpen) {
        let filterInnerContainer = element.querySelector(".filter-inner-container");
        let filterArrowUp = element.querySelector(".filter-arrow-up");
        let filterArrowDown = element.querySelector(".filter-arrow-down");

        filterInnerContainer.classList.toggle("hide", !isOpen);

        filterArrowUp.setAttribute("aria-expanded", isOpen.toString());
        filterArrowDown.setAttribute("aria-expanded", (!isOpen).toString());

        filterArrowUp.classList.toggle("hide", !isOpen);
        filterArrowDown.classList.toggle("hide", isOpen);
    }

    applyFilters() {
        let url = this.clearFilters.getAttribute("href");
        this.filters.forEach(function (filter) {
            if (!filter.value) {
                return
            }
            if (filter.checked || filter.type !== "checkbox") {
                url += "&" + filter.name + "=" + filter.value
            }
        });
        window.location.href = url;
    }

    toggleFilter(event) {
        const parent = event.target.closest('.moj-filter__options');
        const innerContainer= parent.querySelector(
            ".js-options-container"
        );
        innerContainer.classList.toggle("hide");
    }
}
