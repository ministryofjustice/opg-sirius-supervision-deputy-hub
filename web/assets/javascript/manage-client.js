export default class ManageClient {
    constructor(element) {
        this.data = {
            selectedTasks: 0,
        };

        this.checkBoxElements = element.querySelectorAll(".js-mt-checkbox");
        this.allcheckBoxElements = element.querySelectorAll(".js-mt-checkbox-select-all");
        this.bulkAssuranceVisitTaskButton = element.querySelectorAll(".js-mt-edit-btn");
        this.cancelBulkAssuranceVisitTaskButton = element.querySelectorAll(".js-mt-cancel");
        this.xsrfToken = element.querySelector(".js-xsrfToken");
        this.selectedCountElement = element.querySelectorAll(".js-mt-count");
        this.editPanelDiv = element.querySelectorAll(".js-mt-edit-panel");
        this.baseUrl = document.querySelector("[name=api-base-uri]").getAttribute("content");

        this._setupEventListeners();
        const dueDateInput = document.getElementById("dueDate");
        if (dueDateInput) {
            const today = new Date().toISOString().split("T")[0];
            dueDateInput.value = today;
            dueDateInput.min = today;
        }
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

        this.bulkAssuranceVisitTaskButton.forEach((element) => {
            this._showEditTasksPanel = this._showEditTasksPanel.bind(this);
            element.addEventListener("click", this._showEditTasksPanel);
        });

        this.cancelBulkAssuranceVisitTaskButton.forEach((element) => {
            this._hideEditTasksPanel = this._hideEditTasksPanel.bind(this);
            element.addEventListener("click", this._hideEditTasksPanel);
        });
    }

    _updateDomElements() {
        this.selectedCountElement.forEach((element) => {
            element.innerText = this.data.selectedTasks.toString();
        });
        this.bulkAssuranceVisitTaskButton[0].classList.toggle("hide", this.data.selectedTasks === 0);
    }

    _updateSelectedRowStyles(element) {
        element.parentElement.parentElement.parentElement.classList.toggle("govuk-table__select", element.checked);
        element.parentElement.parentElement.parentElement.parentElement.classList.toggle("selected", element.checked);
    }

    _updateSelectedState(event) {
        event.target.checked ? this.data.selectedTasks++ : this.data.selectedTasks--;
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

    _showEditTasksPanel() {
        this.editPanelDiv.forEach((element) => {
            element.classList.toggle("hide", this.data.selectedTasks === 0);
        });
    }

    _hideEditTasksPanel() {
        this.editPanelDiv.forEach((element) => {
            element.classList.toggle("hide", true);
        });
    }
}
