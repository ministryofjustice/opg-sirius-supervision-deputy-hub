export default class DownloadChecker {
    constructor(element) {
        this.links = element.querySelectorAll(".document-download-link");
        this.setupEventListeners();
    }

    setupEventListeners() {
        this.links.forEach((link) => {
            link.addEventListener("click", async function (e) {
                e.preventDefault();
                const deputyId = this.getAttribute("data-deputy-id");
                const docId = this.getAttribute("data-document-id");
                const checkUrl = `/${deputyId}/documents/${docId}/check`;
                const downloadUrl = this.getAttribute("href"); // Get the original download URL

                // Hide error banner in case it's visible
                const banner = document.getElementById("error-banner");
                const infectedLabel = document.getElementById("infected-label-"+docId);
                if (banner) {
                    banner.hidden = true;
                }

                try {
                    // HEAD request checks file availability and whether its infected
                    const headResponse = await fetch(checkUrl, {
                        method: "HEAD",
                        credentials: "include",
                    });

                    if (!headResponse.ok) {
                        throw new Error("File is infected or unavailable");
                    }

                    // Trigger download by navigating to the original download URL
                    window.location.href = downloadUrl;
                } catch (err) {
                    if (banner) {
                        banner.hidden = false;
                        infectedLabel.hidden = false;
                        document.getElementById("error-banner-message").textContent =
                            "This file is blocked. A suspected virus has been detected. Please request a different file from the sender and notify the Implementation team";
                    }
                }
            });
        });
    }
}
