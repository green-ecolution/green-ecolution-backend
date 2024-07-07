import { createFileRoute } from "@tanstack/react-router";
import { Separator } from "../components/ui/separator";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../components/ui/table";
import { Button } from "../components/ui/button";
import { Edit, PlusCircleIcon, Trash, Filter } from "lucide-react";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "../components/ui/popover";
import { Checkbox } from "../components/ui/checkbox";

export const Route = createFileRoute("/vehicles")({
  component: Vehicles,
});

const vehicles = [
  {
    title: "LF 10",
    licensePlate: "FL-TB-1235",
    type: "Wasserfahrzeug",
    location: "Klärwerk",
    status: "Verfügbar",
  },
  {
    title: "LF 20",
    licensePlate: "FL-TB-1235",
    type: "Wasserfahrzeug",
    location: "TBZ Standort",
    status: "Verfügbar",
  },
  {
    title: "LF 10",
    licensePlate: "FL-TB-1235",
    type: "Wasserfahrzeug",
    location: "TBZ Standort",
    status: "Verfügbar",
  },
  {
    title: "LF 20",
    licensePlate: "FL-TB-1235",
    type: "Pritschenwagen",
    location: "Klärwerk",
    status: "Nicht verfügbar",
  },
];

function Vehicles() {
  return (
    <div>
      <div className="h-[48px] flex items-center justify-between mx-4">
        <h1 className="font-bold text-xl">Fahrzeuge</h1>

        <div className="flex items-center gap-2">
        <Popover>
            <PopoverTrigger asChild>
              <Button variant="ghost" size="icon">
                <Filter className="w-6 h-6" />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-80 flex flex-col gap-3">
              <h1 className="font-bold text-xl">Filter</h1>
              <div>
                <h2 className="font-bold ml-1">Status</h2>
                <ul className="ml-2">
                  <li>
                    <Checkbox id="statusVerfuegbar" />
                    <label htmlFor="statusVerfuegbar" className="ml-1">
                      Verfügbar
                    </label>
                  </li>
                  <li>
                    <Checkbox id="statusNichtVerfuegbar" />
                    <label htmlFor="statusNichtVerfuegbar" className="ml-1">
                      Nicht Verfügbar
                    </label>
                  </li>
                </ul>
              </div>
              <div>
                <h2 className="font-bold ml-1">Typ</h2>
                <ul className="ml-2">
                  <li>
                    <Checkbox id="gaertner" />
                    <label htmlFor="gaertner" className="ml-1">
                      Wasserfahrzeug
                    </label>
                  </li>
                  <li>
                    <Checkbox id="foerster" />
                    <label htmlFor="foerster" className="ml-1">
                      Pritschenwagen
                    </label>
                  </li>
                </ul>
              </div>
              <div>
                <h2 className="font-bold ml-1">Standort</h2>
                <ul className="ml-2">
                  <li>
                    <Checkbox id="gaertner" />
                    <label htmlFor="gaertner" className="ml-1">
                      Klärwerk
                    </label>
                  </li>
                  <li>
                    <Checkbox id="foerster" />
                    <label htmlFor="foerster" className="ml-1">
                      TBZ
                    </label>
                  </li>
                </ul>
              </div>
            </PopoverContent>
          </Popover>
          <Button variant="default">
            <PlusCircleIcon className="w-4 h-4" />
            <span className="ml-2">Fahrzeug hinzufügen</span>
          </Button>
        </div>
      </div>
      <Separator />

      <div className="p-4">
        <div className="flex justify-end items-center"></div>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[100px]">Beziechnung</TableHead>
              <TableHead>Kennzeichen</TableHead>
              <TableHead>Typ</TableHead>
              <TableHead>Standort</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Aktion</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {vehicles.map((vehicle) => (
              <TableRow key={vehicle.title}>
                <TableCell className="font-medium">{vehicle.title}</TableCell>
                <TableCell>{vehicle.licensePlate}</TableCell>
                <TableCell>{vehicle.type}</TableCell>
                <TableCell>{vehicle.location}</TableCell>
                <TableCell>{vehicle.status}</TableCell>
                <TableCell className="text-right">
                  <Button variant="ghost" size="icon">
                    <Edit className="w-4 h-4" />
                  </Button>
                  <Button variant="ghost" size="icon">
                    <Trash className="w-4 h-4" />
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
