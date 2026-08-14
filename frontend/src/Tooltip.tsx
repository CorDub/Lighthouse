import { Portal } from "solid-js/web";
import { onMount, createSignal, mergeProps } from "solid-js";
import "./styles/Tooltip.css";

type TooltipPosition = "over" | "right" | "left" | "under";

type TooltipProps = {
  position: TooltipPosition,
  topSender: number,
  leftSender: number,
  widthSender: number,
  heightSender: number,
  text: string,
  topSpaceInRem?: number,
  leftSpaceInRem?: number,
}

function Tooltip(tooltipProps: TooltipProps) {
  let tooltipContentRef: HTMLDivElement | undefined;
  let arrowRef: HTMLDivElement | undefined;
  const [tooltipRect, setTooltipRect] = createSignal<DOMRect>()
  const [arrowRect, setArrowRect] = createSignal<DOMRect>()
  
  const props = mergeProps(
    { 
      topSpaceInRem : 0.5,
      leftSpaceInRem : 0.5
    },
    tooltipProps
  );

  onMount(() => {
    setTooltipRect(tooltipContentRef?.getBoundingClientRect());
    setArrowRect(arrowRef?.getBoundingClientRect())
  })

  function getTopPosition() {
    // get the tooltip height with the text passed as props;
    const rect = tooltipRect();
    if (!rect) return;

    const tooltipHeight = rect.height

    // get the space between the tooltip and tooltiped element;
    const rootFontSizePx = parseFloat(getComputedStyle(document.documentElement).fontSize);
    const topSpaceInPixels = props.topSpaceInRem *  rootFontSizePx;

    //get the halves for the left and right cases;
    const middleElement = props.heightSender / 2;
    const middleTooltip = tooltipHeight / 2;

    const halvesDiff = middleTooltip - middleElement;
    
    // return final top position
    let finalTop: number;
    switch (props.position) {
      case "over": 
        finalTop = props.topSender - tooltipHeight - topSpaceInPixels;
        break;
      case "right":
        finalTop = props.topSender - halvesDiff;
        break;
      case "left":
        finalTop = props.topSender - halvesDiff;
        break;
      case "under":
        finalTop = props.topSender + props.heightSender + topSpaceInPixels;
        break;
      default: 
        finalTop = props.topSender;
    }
    return finalTop;
  }

  function getLeftPosition() {
    // get the tooltip width with the text passed as props;
    const rect = tooltipRect();
    if (!rect) return;

    const tooltipWidth = rect.width

    //get the element width;
    const elementWidth = props.widthSender;

    //get the middle of both elements to center them;
    const middleElement = elementWidth / 2;
    const middleTooltip = tooltipWidth / 2;

    //get the difference between each halves
    const halvesDiff = middleTooltip - middleElement;

    //get the extra space
    const rootFontSizePx = parseFloat(getComputedStyle(document.documentElement).fontSize);
    const leftSpaceInPixels = props.leftSpaceInRem * rootFontSizePx;

    //return the left position with median points aligned;
    let finalLeft: number;
    switch (props.position) {
      case "over":
        finalLeft = props.leftSender - halvesDiff
        break;
      case "right":
        finalLeft = props.leftSender + props.widthSender + leftSpaceInPixels;
        break;
      case "left":
        finalLeft = props.leftSender - leftSpaceInPixels - tooltipWidth;
        break;
      case "under":
        finalLeft = props.leftSender - halvesDiff
        break;
      default: 
        finalLeft = props.leftSender - halvesDiff
    }

    return finalLeft;
  }

  function getTopArrowPosition() {
    let finalArrowTop: number;

    const rect = arrowRect();
    const ttrect = tooltipRect();
    if (!rect || ! ttrect) return;

    const arrowHeight = rect.height
    const arrowHalf = arrowHeight / 2;

    switch (props.position) {
      case "over":
        finalArrowTop = ttrect.bottom - arrowHalf - 3;
        break;
      case "right":
        finalArrowTop = ttrect.top + (ttrect.height / 2) - 7.5;
        break;
      case "left":
        // finalArrowTop = props.topSender - halvesDiff
        finalArrowTop = ttrect.top + (ttrect.height / 2) - 7.5;
        break;
      case "under":
        // finalArrowTop = props.topSender - props.heightSender + arrowHalf
        finalArrowTop = ttrect.top - arrowHalf + 9;
        break;
      default: 
        finalArrowTop = props.topSender
    }
    
    return finalArrowTop;
  }

  function getLeftArrowPosition() {
    let leftArrowTop: number;

    const rect = arrowRect();
    const ttrect = tooltipRect();
    if (!rect || ! ttrect) return;

    switch (props.position) {
      case "over":
        leftArrowTop = (ttrect.width / 2) - 7.5;
        break;
      case "right":
        // leftArrowTop = props.leftSender + props.widthSender - (arrowWidth / 2)
        leftArrowTop = ttrect.left - 3;
        break;
      case "left":
        // leftArrowTop = props.leftSender + (arrowWidth / 2)
        leftArrowTop = ttrect.right - 13;
        break;
      case "under":
        // leftArrowTop = props.leftSender + (props.widthSender / 2) - (arrowWidth / 2)
        leftArrowTop = (ttrect.width / 2) - 7.5;
        break;
      default: 
        leftArrowTop = props.leftSender
    }
      
    return leftArrowTop;
  }

  return(
    <Portal>
      <div class="tooltip"
        style={{
          "top": `${getTopPosition() ?? 0}px`,
          "left": `${getLeftPosition() ?? 0}px`,
        }}>
        <div class="tooltip-content"
          ref={tooltipContentRef}>
          <p>{props.text}</p>
          <div class="tooltip-arrow"
            ref={arrowRef}
            style={{
              "top": `${getTopArrowPosition() ?? 0}px`,
              "left": `${getLeftArrowPosition() ?? 0}px`
            }}></div>
        </div>
      </div>
    </Portal>
  )
}

export default Tooltip;