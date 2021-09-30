const MojBannerAutoHide = (element) => {
  
    const elements = element.querySelectorAll(
      '[data-module="moj-banner-auto-hide"]'
    );
  
    elements.forEach((bannerElement) => {
      setTimeout(function() {
        bannerElement.classList.add('hide');
      }, 1000)
    });
  };
  
  export default MojBannerAutoHide;
  