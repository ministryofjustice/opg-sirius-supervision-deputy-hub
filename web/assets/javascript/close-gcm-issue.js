export default class CloseGcmIssue {
    constructor(element) {
        this.data = {
            selectedTasks: 0,
        };

        this.checkBoxElements = element.querySelectorAll(".js-mt-checkbox");
        this.allcheckBoxElements = element.querySelectorAll(
            ".js-mt-checkbox-select-all"
        );
        this.gcmClosedIssueButton = element.querySelectorAll(".js-mt-edit-btn");
        this.xsrfToken = element.querySelector(".js-xsrfToken");

        this._setupEventListeners();
    }

    _setupEventListeners() {
        this.checkBoxElements.forEach((element) => {
            this._updateSelectedState = this._updateSelectedState.bind(this);
            element.addEventListener("click", this._updateSelectedState);
        });

        this.allcheckBoxElements.forEach((element) => {
            this._updateAllSelectedState = this._updateAllSelectedState.bind(this);
            element.addEventListener("click", this._updateAllSelectedState);
        });
    }

    _updateDomElements() {
        this.gcmClosedIssueButton[0].classList.toggle(
            "hide",
            this.data.selectedTasks === 0
        );
    }

    _updateSelectedRowStyles(element) {
        element.parentElement.parentElement.parentElement.classList.toggle(
            "govuk-table__select",
            element.checked
        );
        element.parentElement.parentElement.parentElement.parentElement.classList.toggle(
            "selected",
            element.checked
        );
    }

    _updateSelectedState(event) {
        event.target.checked
            ? this.data.selectedTasks++
            : this.data.selectedTasks--;
        this._updateSelectedRowStyles(event.target);
        this._updateDomElements();
    }

    _updateAllSelectedState(event) {
        let isChecked = event.target.checked;

        this.checkBoxElements.forEach((checkbox) => {
            checkbox.checked = isChecked;

            this._updateSelectedRowStyles(checkbox);
        });

        this.data.selectedTasks = isChecked ? this.checkBoxElements.length : 0;
        this._updateDomElements();
    }
}
