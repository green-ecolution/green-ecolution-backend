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

export const Route = createFileRoute("/team")({
  component: Team,
});

enum Status {
  Verfuegbar = "Verfügbar",
  NichtVerfuegbar = "Nicht Verfügbar",
}

const team = [
  {
    name: "Hans Olaf",
    jobPosition: "Gärtner",
    status: Status.Verfuegbar,
  },
  {
    name: "Timo Müller",
    jobPosition: "Förster",
    status: Status.Verfuegbar,
  },
  {
    name: "Dieter Jürgensen",
    jobPosition: "Gärtner",
    status: Status.NichtVerfuegbar,
  },
  {
    name: "Ralf Peter",
    jobPosition: "Gärtner",
    status: Status.Verfuegbar,
  },
  {
    name: "Harald Thomsen",
    jobPosition: "Förster",
    status: Status.Verfuegbar,
  },
  {
    name: "Uwe Schmidt",
    jobPosition: "Förster",
    status: Status.NichtVerfuegbar,
  },
];

function Team() {
  return (
    <div>
      <div className="h-[48px] flex items-center justify-between mx-4">
        <h1 className="font-bold text-xl">Team Mitglieder</h1>

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
                <h2 className="font-bold ml-1">Job Position</h2>
                <ul className="ml-2">
                  <li>
                    <Checkbox id="gaertner" />
                    <label htmlFor="gaertner" className="ml-1">
                      Gärtner
                    </label>
                  </li>
                  <li>
                    <Checkbox id="foerster" />
                    <label htmlFor="foerster" className="ml-1">
                      Förster
                    </label>
                  </li>
                </ul>
              </div>
            </PopoverContent>
          </Popover>
          <Button variant="default">
            <PlusCircleIcon className="w-4 h-4" />
            <span className="ml-2">Mitglied hinzufügen</span>
          </Button>
        </div>
      </div>
      <Separator />

      <div className="p-4">
        <div className="flex justify-end items-center"></div>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Job Position</TableHead>
              <TableHead>Aktueller Status</TableHead>
              <TableHead className="text-right">Aktion</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {team.map((teamMember) => (
              <TableRow key={teamMember.name}>
                <TableCell className="font-medium">{teamMember.name}</TableCell>
                <TableCell>{teamMember.jobPosition}</TableCell>
                <TableCell>{teamMember.status}</TableCell>
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
