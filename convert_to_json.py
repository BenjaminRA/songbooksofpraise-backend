#!/usr/bin/env python3
import re
import json

def parse_sql_to_json(sql_file_path, output_file_path):
    with open(sql_file_path, 'r', encoding='utf-8') as file:
        content = file.read()
    
    # Dictionary to store countries and their states with IDs
    countries_data = {}
    country_id_counter = 1
    
    # Parse countries first
    countries_pattern = r"INSERT INTO countries \(name, iso_alpha2, iso_alpha3, iso_numeric\) VALUES(.*?);"
    countries_match = re.search(countries_pattern, content, re.DOTALL)
    
    if countries_match:
        countries_values = countries_match.group(1)
        # Extract country entries with full details
        country_entries = re.findall(r"\('([^']+)',\s*'([^']+)',\s*'([^']+)',\s*'([^']+)'\)", countries_values)
        
        # Initialize countries dictionary with ISO codes as keys and assign IDs
        for name, iso_alpha2, iso_alpha3, iso_numeric in country_entries:
            countries_data[iso_alpha3] = {
                "id": country_id_counter,
                "name": name,
                "iso_alpha2": iso_alpha2,
                "iso_alpha3": iso_alpha3,
                "iso_numeric": iso_numeric,
                "states": []
            }
            country_id_counter += 1
    
    # Parse states
    states_pattern = r"INSERT INTO states \(name, country_id\) VALUES(.*?)COMMIT;"
    states_match = re.search(states_pattern, content, re.DOTALL)
    
    state_id_counter = 1
    if states_match:
        states_values = states_match.group(1)
        # Extract state entries
        state_entries = re.findall(r"\('([^']+(?:''[^']*)*)',\s*\(SELECT id FROM countries WHERE iso_alpha3 = '([^']+)'\)\)", states_values)
        
        # Add states to their respective countries with IDs
        for state_name, iso_code in state_entries:
            # Handle escaped single quotes
            state_name = state_name.replace("''", "'")
            
            if iso_code in countries_data:
                state_data = {
                    "id": state_id_counter,
                    "name": state_name,
                    # "country_id": countries_data[iso_code]["id"]
                }
                countries_data[iso_code]["states"].append(state_data)
                state_id_counter += 1
    
    # Convert to the desired JSON structure (array of countries with states)
    result = []
    for iso_code, country_data in countries_data.items():
        # Sort states alphabetically by name
        sorted_states = sorted(country_data["states"], key=lambda x: x["id"])
        
        result.append({
            "id": country_data["id"],
            "name": country_data["name"],
            "iso_alpha2": country_data["iso_alpha2"],
            "iso_alpha3": country_data["iso_alpha3"],
            "iso_numeric": country_data["iso_numeric"],
            "states": sorted_states
        })
    
    # Sort countries alphabetically by name
    result.sort(key=lambda x: x["name"])

    # Minify JSON output
    minified_result = json.dumps(result, ensure_ascii=False)

    # Write to JSON file
    with open(output_file_path, 'w', encoding='utf-8') as output_file:
        output_file.write(minified_result)

    print(f"Successfully converted SQL data to JSON!")
    print(f"Found {len(result)} countries")
    print(f"Total states/divisions: {sum(len(country['states']) for country in result)}")
    print(f"Output saved to: {output_file_path}")
    
    return result

if __name__ == "__main__":
    sql_file = "churches.sql"
    json_file = "countries_states.json"
    
    try:
        data = parse_sql_to_json(sql_file, json_file)
        
        # # Print some sample data
        # print("\nSample data:")
        # for country in data[:3]:  # Show first 3 countries
        #     print(f"\n{country['name']} ({country['iso_alpha3']}) - ID: {country['id']}")
        #     print(f"  {len(country['states'])} states/divisions")
        #     if country['states']:
        #         print(f"  First few states:")
        #         for state in country['states'][:3]:
        #             print(f"    - {state['name']} (ID: {state['id']})")
        #         if len(country['states']) > 3:
        #             print("    ...")
    
    except Exception as e:
        print(f"Error: {e}")
