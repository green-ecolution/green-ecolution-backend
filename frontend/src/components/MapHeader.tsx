import { Filter, Menu } from "lucide-react";
import { Button } from "./ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "./ui/popover";
import { Checkbox } from "./ui/checkbox";
import { Input } from "./ui/input";

export interface HeaderProps {}

export interface Filter {
  statusHealthy: boolean;
  statusNeutral: boolean;
  statusUnhealthy: boolean;
  oneYear: boolean;
  twoYear: boolean;
  threeYear: boolean;
}

const MapHeader = ({}: HeaderProps) => {
  return (
    <div className="z-50 absolute top-4 left-4 w-[450px] h-12 bg-white rounded shadow-lg">
      <div className="flex justify-between items-center h-full mx-2 gap-2">
        <div className="flex items-center">
          <Button variant="ghost" size="icon">
            <Menu className="w-6 h-6" />
          </Button>
        </div>
        <Input type="email" placeholder="Suche nach Baum..." />
        <div>
          <div className="flex items-center gap-1">
            <div className="flex items-center">
              <Popover>
                <PopoverTrigger asChild>
                  <Button variant="ghost" size="icon">
                    <Filter className="w-6 h-6" />
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-80 flex flex-col gap-3">
                  <h1 className="font-bold text-xl">Filter</h1>
                  <div>
                    <h2 className="font-bold ml-1">Tree Status</h2>
                    <ul className="ml-2">
                      <li>
                        <Checkbox id="statusHealthy" />
                        <label htmlFor="statusHealthy" className="ml-1">
                          Healthy
                        </label>
                      </li>
                      <li>
                        <Checkbox id="statusNeutral" />
                        <label htmlFor="statusNeutral" className="ml-1">
                          Neutral
                        </label>
                      </li>
                      <li>
                        <Checkbox id="statusUnhealthy" />
                        <label htmlFor="statusUnhealthy" className="ml-1">
                          Unhealthy
                        </label>
                      </li>
                    </ul>
                  </div>
                  <div>
                    <h2 className="font-bold ml-1">Tree Age</h2>
                    <ul className="ml-2">
                      <li>
                        <Checkbox id="oneYear" />
                        <label htmlFor="oneYear" className="ml-1">
                          1 year
                        </label>
                      </li>
                      <li>
                        <Checkbox id="twoYear" />
                        <label htmlFor="twoYear" className="ml-1">
                          2 years
                        </label>
                      </li>
                      <li>
                        <Checkbox id="threeYear" />
                        <label htmlFor="threeYear" className="ml-1">
                          3 years
                        </label>
                      </li>
                    </ul>
                  </div>
                </PopoverContent>
              </Popover>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default MapHeader;
