import { For, createSignal } from "solid-js";
import "./styles/SearchResults.css";

type SearchResultsProps = {
  searchResults: string[],
  resultSelectFn: (...args: any[]) => any,
  arg?: any
}

function SearchResults(props: SearchResultsProps) {
  const [activeIndex, setActiveIndex] = createSignal(0)

  return (
    <div class="search-results">
      <For each={props.searchResults}>
        {(res, i) => 
          <div class="sr-line"
            classList={{srlactive: activeIndex() === i()}}
            onMouseDown={(e) => {
              e.preventDefault();
              props.resultSelectFn([...props.arg, {name: res, inDB: true}])
            }}
            onMouseEnter={() => setActiveIndex(i())}
            onMouseLeave={() => setActiveIndex(0)}>
            {res}
          </div>
        }
      </For>
    </div>
  )
}

export default SearchResults;